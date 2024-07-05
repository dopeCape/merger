package broker

import (
	"github.com/dopeCape/schduler/pkg/shared"
	"github.com/hibiken/asynq"
)

var broker *asynq.Client

func RunBroker(config shared.Config) (*Brokers, *asynq.Client) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisAddress})
	brokers := &Brokers{client}
	broker = client

	return brokers, client
}

func GetBroker() *asynq.Client {
	return broker
}
