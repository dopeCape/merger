package broker

import (
	"encoding/json"
	"time"

	"github.com/dopeCape/schduler/internal/models"
	"github.com/hibiken/asynq"
)

type Brokers struct {
	client *asynq.Client
}

func (br *Brokers) Enqueue(url string, body json.RawMessage, headers map[string]string) (*asynq.TaskInfo, error) {
	task, err := NewPostTask(url, body, headers)
	if err != nil {
		return nil, err
	}
	info, err := br.client.Enqueue(task)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (br *Brokers) EnqueueWithDelay(url string, body json.RawMessage, headers map[string]string, delay time.Duration) (*asynq.TaskInfo, error) {
	task, err := NewPostTask(url, body, headers)
	if err != nil {
		return nil, err
	}
	info, err := br.client.Enqueue(task, asynq.ProcessIn(time.Second*delay))
	if err != nil {
		return nil, err
	}
	return info, nil

}

func (br *Brokers) EnqueueAt(url string, body json.RawMessage, headers map[string]string, timeToRun string) (*asynq.TaskInfo, error) {
	task, err := NewPostTask(url, body, headers)
	if err != nil {
		return nil, err
	}
	processAt, err := time.Parse(time.RFC3339Nano, timeToRun)
	info, err := br.client.Enqueue(task, asynq.ProcessAt(processAt))
	if err != nil {
		return nil, err
	}
	return info, nil

}

type PostDeliveryPayload struct {
	URL     string
	Body    json.RawMessage
	Headers map[string]string
}

func NewPostTask(URL string, body json.RawMessage, headers map[string]string) (*asynq.Task, error) {
	payload, err := json.Marshal(PostDeliveryPayload{URL: URL, Body: body, Headers: headers})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask("post:task", payload, asynq.MaxRetry(8), asynq.Timeout(time.Minute*15)), nil
}

func NewTaskUpdateJob(task models.Task) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask("task:update", payload, asynq.MaxRetry(1)), nil
}

func EnqueueNewTaskUpdateJob(task models.Task) (*asynq.TaskInfo, error) {
	br := GetBroker()
	createdTask, err := NewTaskUpdateJob(task)
	if err != nil {
		return nil, err
	}
	info, err := br.Enqueue(createdTask, asynq.ProcessIn(time.Second*2))
	if err != nil {
		return nil, err
	}
	return info, nil

}
