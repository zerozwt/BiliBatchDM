package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync/atomic"

	jsoniter "github.com/json-iterator/go"
)

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter

	Next func() error

	once *int32
}

type Handler func(ctx *Context) error

func HandlerChain(handlers ...Handler) http.HandlerFunc {
	if len(handlers) == 0 {
		handlers = append(handlers, func(ctx *Context) error {
			ctx.WriteResponse(404, "Not found", nil)
			return nil
		})
	}

	ret := func(rsp http.ResponseWriter, req *http.Request) {
		chain := make([]*Context, len(handlers))
		var write_lock int32 = 0
		for idx := range handlers {
			chain[idx] = &Context{
				Request:  req,
				Response: rsp,
				once:     &write_lock,
			}
		}

		total := len(chain)
		chain = append(chain, nil)
		handlers = append(handlers, nil)

		for i := 0; i < total; i++ {
			next := handlers[i+1]
			ctx := chain[i+1]
			chain[i].Next = func() error {
				return next(ctx)
			}
		}

		handlers[0](chain[0])
	}

	return ret
}

func (ctx *Context) ReadRequest(req interface{}) error {
	tag := "form"
	value := reflect.ValueOf(req).Elem()
	typ := value.Type()

	if err := ctx.Request.ParseForm(); err != nil {
		logger().Printf("HTTP request parse form failed: %v", err)
		return err
	}

	for i := 0; i < typ.NumField(); i++ {
		field_t := typ.Field(i)
		key, ok := field_t.Tag.Lookup(tag)
		if !ok || len(key) == 0 {
			continue
		}

		if _, ok := ctx.Request.Form[key]; !ok {
			err := fmt.Errorf("param '%s' is needed but not provided", key)
			logger().Printf("HTTP request parse form failed: %v", err)
			return err
		}
		param_value := ctx.Request.Form.Get(key)

		switch value.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(param_value, 10, 64)
			if err != nil {
				logger().Printf("HTTP request parse form failed: %v", err)
				return err
			}
			value.Field(i).SetInt(n)
		case reflect.String:
			value.Field(i).SetString(param_value)
		default:
			err := fmt.Errorf("type of param %s is not supported: %T", key, value.Field(i).Interface())
			logger().Printf("HTTP request parse form failed: %v", err)
			return err
		}
	}
	return nil
}

func (ctx *Context) ReadRequestJson(req interface{}) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(ctx.Request.Body)
	return decoder.Decode(req)
}

func (ctx *Context) WriteResponse(code int, msg string, data interface{}) {
	if !atomic.CompareAndSwapInt32(ctx.once, 0, 1) {
		return
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	tmp := map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}
	tmp_data, _ := json.Marshal(tmp)
	ctx.Response.Header().Set("Content-Type", "application/json")
	if _, err := ctx.Response.Write(tmp_data); err != nil {
		logger().Printf("HTTP write response failed: %v", err)
	}
}
