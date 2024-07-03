package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hibiken/asynq"
	"github.com/tinkler/keeper_grid/internal/handle/automation_compatible"
	"github.com/tinkler/keeper_grid/pkg/tasks"
	"github.com/tinkler/moonmist/pkg/mlog"
)

func main() {
	mlog.ConsoleLevel = mlog.L_DEBUG
	root, _ := os.Getwd()
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// progress
	wallet := keystore.NewKeyStore(filepath.Join(root, "wallet"), keystore.StandardScryptN, keystore.StandardScryptP)
	_, err := wallet.NewAccount("nsdcawer512")
	if err != nil {
		panic(err)
	}
	iotexPool, err := automation_compatible.NewHandler(serverCtx, "http://localhost:8545", wallet)
	if err != nil {
		panic(err)
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379", Password: "password"},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.ACT_IOTEX, iotexPool.HandleCheckUpkeep)

	signalCtx, signalStop := signal.NotifyContext(serverCtx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signalStop()

	go func() {
		if err := srv.Run(mux); err != nil {
			panic(err)
		}
	}()

	go func() {
		<-signalCtx.Done()
		// shutdownCtx, cancelShutdownCtx := context.WithTimeout(serverCtx, 30*time.Second)
		// defer cancelShutdownCtx()

		// go func() {
		// 	<-shutdownCtx.Done()
		// 	if shutdownCtx.Err() == context.DeadlineExceeded {
		// 		mlog.Error("graceful shutdown time out... forcing exit.")
		// 	}
		// }()

		srv.Shutdown()
		serverStopCtx()
	}()

	<-serverCtx.Done()
}
