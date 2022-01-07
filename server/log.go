package main

import (
	"io"
	"log"
	"os"
	"sync/atomic"
)

var g_logger atomic.Value

func init() {
	SetLogWriter(os.Stdout)
}

func logger() *log.Logger {
	return g_logger.Load().(*log.Logger)
}

func SetLogWriter(out io.Writer) {
	if out == nil {
		out = os.Stdout
	}
	logger := log.New(out, "BBDMA ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	g_logger.Store(logger)
}
