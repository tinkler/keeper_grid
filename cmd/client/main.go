package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/tinkler/keeper_grid/pkg/tasks/automation_compatible"
	"github.com/tinkler/moonmist/pkg/mlog"
)

func main() {
	mlog.ConsoleLevel = mlog.L_DEBUG
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379", Password: "password"})
	task, err := automation_compatible.NewTask("0x13AE1e5702DffcfD9900713f2Ba9d2995af7ed5D", "0xaf64fa33a8a0640e1eaf6b6f32b81efd11653460")
	if err != nil {
		panic(err)
	}

	signalCtx, signalStop := signal.NotifyContext(serverCtx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signalStop()

	go func() {
		for {

			info, err := client.Enqueue(task, asynq.Timeout(3*time.Minute))
			if err != nil {
				panic(err)
			}
			mlog.Info("enqueued task: id=%s queue=%s", info.ID, info.Queue)
			time.Sleep(time.Minute * 1)
		}
	}()

	go func() {
		<-signalCtx.Done()

		client.Close()
		serverStopCtx()
	}()

	<-serverCtx.Done()
}
