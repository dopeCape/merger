package inspector

import (
	"errors"

	"github.com/hibiken/asynq"
)

func (in *Inspector) Dequque(queue string, id string) error {
	err := in.inspector.DeleteTask(queue, id)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskNotFound) {
			return errors.New("Task not found")
		}
		if errors.Is(err, asynq.ErrQueueNotEmpty) {
			return errors.New("queue not found")
		}
		return err
	}
	return nil
}
