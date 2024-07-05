package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
	"github.com/dopeCape/schduler/pkg/suscriber"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func GetSchduler() *asynq.PeriodicTaskManager {
	provider := &RDBBasedConfigProvider{}

	mgr, err := asynq.NewPeriodicTaskManager(
		asynq.PeriodicTaskManagerOpts{
			RedisConnOpt:               asynq.RedisClientOpt{Addr: "localhost:6379"},
			PeriodicTaskConfigProvider: provider,         // this provider object is the interface to your config source
			SyncInterval:               30 * time.Second, // this field specifies how often sync should happen
		})
	if err != nil {
		log.Fatal(err)
	}

	return mgr
}

// FileBasedConfigProvider implements asynq.PeriodicTaskConfigProvider interface.
type RDBBasedConfigProvider struct {
}

// Parses the yaml file and return a list of PeriodicTaskConfigs.
func (p *RDBBasedConfigProvider) GetConfigs() ([]*asynq.PeriodicTaskConfig, error) {
	db, err := rdb.GetDb()
	if err != nil {
		panic("Failed to get db in schduler")
	}
	var tasks []models.Task
	res := db.Where("is_cron = ?", true).Find(&tasks)
	if res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			fmt.Print(res.Error)
			panic("Failed to retive cron jobs")
		}
	}

	var configs []*asynq.PeriodicTaskConfig
	for _, task := range tasks {
		payloadByte, _ := json.Marshal(suscriber.CronPayload{Body: json.RawMessage(task.Payload), Headers: task.Headers, URL: task.URL})
		configs = append(configs, &asynq.PeriodicTaskConfig{Cronspec: task.CronExpresion, Task: asynq.NewTask("cron:task", payloadByte, asynq.TaskID(task.ID), asynq.MaxRetry(8))})
	}
	return configs, nil
}

type Config struct {
	Cronspec string
	TaskType string
	Payload  json.RawMessage
}
