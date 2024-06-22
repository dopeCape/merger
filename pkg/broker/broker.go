package broker

import (
	"github.com/dopeCape/schduler/pkg/shared"
	"github.com/hibiken/asynq"
)

func RunBroker(config shared.Config) (*Brokers, *asynq.Client) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisAddress})
	brokers := &Brokers{client}
	return brokers, client
}
