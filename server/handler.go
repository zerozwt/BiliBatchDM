package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Err233 error = errors.New("233")

// general handlers

func AccessLog(ctx *Context) error {
	req := ctx.Request
	now := time.Now()
	url := req.URL.String()
	defer func() {
		if url == "/api/progress" {
			return
		}
		logger().Printf("ACCESS %s %s , process time: %v", req.Method, url, time.Since(now))
	}()
	return ctx.Next()
}

type APIHandler func(ctx *Context) (interface{}, error)

func WithLog(handlers ...Handler) http.HandlerFunc {
	handlers = append([]Handler{AccessLog}, handlers...)
	return HandlerChain(handlers...)
}

func API(api APIHandler) http.HandlerFunc {
	wrapper := func(ctx *Context) error {
		rsp, err := api(ctx)
		if err != nil {
			if errors.Is(err, Err233) {
				ctx.WriteResponse(233, "", rsp)
			} else {
				ctx.WriteResponse(500, err.Error(), rsp)
			}
			return err
		}
		ctx.WriteResponse(0, "", rsp)
		return nil
	}
	return WithLog(wrapper)
}

// file handler

func ServeSPA(ctx *Context) error {
	if strings.Contains(ctx.Request.URL.Path, "..") || strings.Contains(ctx.Request.URL.Path, "/.") {
		http.Redirect(ctx.Response, ctx.Request, "/404", http.StatusMovedPermanently)
		return nil
	}

	target := ctx.Request.URL.Path
	if strings.HasSuffix(target, "/") {
		target += "index.html"
	}
	target = filepath.Join(www_root, target)

	if f, err := os.Open(target); err != nil {
		ctx.Request.URL.Path = "/" // SPA support
	} else {
		f.Close()
	}

	http.FileServer(http.Dir(www_root)).ServeHTTP(ctx.Response, ctx.Request)
	return nil
}

// API handlers

func BatchSend(ctx *Context) (interface{}, error) {
	req := TaskRequest{}
	if err := ctx.ReadRequestJson(&req); err != nil {
		return nil, err
	}

	if req.Sender == 0 {
		return nil, fmt.Errorf("invalid sender uid: %d", req.Sender)
	}
	if len(req.Sess) == 0 {
		return nil, fmt.Errorf("empty sess")
	}
	if len(req.JCT) == 0 {
		return nil, fmt.Errorf("empty jct")
	}
	if req.MinSec > req.MaxSec {
		return nil, fmt.Errorf("invalid interval range (%d ~ %d)", req.MinSec, req.MaxSec)
	}
	if len(req.Content) == 0 {
		return nil, fmt.Errorf("empty content")
	}
	if len(req.RecvList) == 0 {
		return nil, fmt.Errorf("empty reciever list")
	}

	return nil, SubmitTask(&req)
}

func QueryProgress(ctx *Context) (interface{}, error) {
	task := CurrentTask.Clone()
	if task == nil {
		return nil, Err233
	}
	return task, nil
}

type PicItem struct {
	Url string `json:"url"`
}

func PicList(ctx *Context) (interface{}, error) {
	list := []PicItem{}
	dir := filepath.Join(www_root, "pics")
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if (strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".png")) && strings.HasPrefix(path, www_root) {
			list = append(list, PicItem{Url: path[len(www_root):]})
		}
		return nil
	})
	return map[string]interface{}{"list": list}, err
}
