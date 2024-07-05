package inspector

import (
	"github.com/dopeCape/schduler/pkg/shared"
	"github.com/hibiken/asynq"
)

type Inspector struct {
	inspector *asynq.Inspector
}

var inspectorSaved *asynq.Inspector

func GetInspector(config shared.Config) (*asynq.Inspector, *Inspector) {
	ins := asynq.NewInspector(asynq.RedisClientOpt{Addr: config.RedisAddress})
	inspectorSaved = ins
	return ins, &Inspector{ins}

}

func GetInspectorSaved() *asynq.Inspector {
	return inspectorSaved
}
