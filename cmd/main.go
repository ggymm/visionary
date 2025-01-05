package main

import (
	"atlas/internal/task"
	"os"
	"os/signal"
	"syscall"
	"time"

	"atlas/internal/api"
	"atlas/pkg/app"
	"atlas/pkg/data"
	"atlas/pkg/log"
)

func init() {
	app.Init()
	log.Init()
	data.Init()
}

func main() {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 启动服务
	go func() {
		err := api.NewServer().Start()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		err := task.NewScanner().Start()
		if err != nil {
			panic(err)
		}
	}()

EXIT:
	for {
		sig := <-sc

		// 信号处理
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	time.Sleep(time.Second)
	os.Exit(state)
}
