package main

import (
	"flag"
	"net"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	dm "github.com/zerozwt/BLiveDanmaku"
)

func main() {
	var port uint
	flag.UintVar(&port, "port", 8008, "http serve port")
	flag.Parse()

	if port <= 0 || port > 65535 {
		logger().Printf("Invalid http listen %d", port)
		return
	}

	// set wwwroot
	path, err := exePath()
	if err != nil {
		logger().Printf("Get current path failed: %v", err)
		return
	}
	www_root = filepath.Join(path, "dist")
	logger().Printf("wwwroot = %s", www_root)

	// get device id
	gDMDevID, err = dm.GetDMDeviceID()
	if err != nil {
		logger().Printf("Get direct message device id failed: %v", err)
		return
	}
	logger().Printf("direct message device id: %s", gDMDevID)

	// init mux
	server := mux.NewRouter()
	api := server.PathPrefix("/api/").Subrouter()
	api.Path("/batch_send").Methods("POST").HandlerFunc(API(BatchSend))
	api.Path("/progress").Methods("GET").HandlerFunc(API(QueryProgress))
	api.PathPrefix("/").HandlerFunc(WithLog(func(ctx *Context) error {
		ctx.WriteResponse(404, "No handlers", nil)
		return nil
	}))
	server.PathPrefix("/").HandlerFunc(WithLog(ServeSPA))

	// start web server
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(port)))
	if err != nil {
		logger().Printf("Listening on localhost:%d failed: %v", port, err)
		return
	}
	logger().Printf("Please open your browser and visit http://localhost:%d/", port)
	http.Serve(lis, server)
}
