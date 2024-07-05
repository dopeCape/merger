package suscriber

import (
	"log"

	"github.com/dopeCape/schduler/pkg/shared"
	"github.com/hibiken/asynq"
)

func Run(config shared.Config) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.RedisAddress},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: config.Concurrency,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("post:task", HandleCallBackTask)
	mux.HandleFunc("cron:task", HandleCronCallBackTask)
	mux.HandleFunc("task:update", HandleSaveTask)
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
