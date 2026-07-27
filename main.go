package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ao9911/go-matrix/log"

	"github.com/ao9911/bluebell-new/conf"
	"github.com/ao9911/bluebell-new/router"
)

func main() {
	flag.Parse()
	// 初始化配置
	conf.Init()
	//初始化日志
	log.Init(conf.Conf.Log)
	router.Init(conf.Conf)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	router.Stop(ctx)
	log.Info("Server Exited Properly")
}
