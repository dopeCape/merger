package suscriber

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	dbactions "github.com/dopeCape/schduler/internal/db_actions"
	"github.com/dopeCape/schduler/internal/models"
	"github.com/hibiken/asynq"
	"github.com/rs/xid"
)

type Payload struct {
	URL     string
	Body    json.RawMessage
	Headers map[string]string
}

type CronPayload struct {
	URL     string
	Body    json.RawMessage
	Headers []string
}

func HandleCallBackTask(ctx context.Context, t *asynq.Task) error {

	fmt.Println("Started post job")
	var p Payload
	taskId := t.ResultWriter().TaskID()
	execution := models.Execution{
		TaskID: taskId,
		Status: models.Active,
		ID:     xid.New().String(),
		RanAt:  time.Now().String(),
	}
	dbactions.CreateExecution(&execution)
	dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Active})
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	bodyReader := bytes.NewReader(p.Body)
	client := &http.Client{}
	req, err := http.NewRequest("POST", p.URL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	for k, v := range p.Headers {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	bodyBytes, _ := io.ReadAll(res.Body)
	now := time.Now().String()
	if res.StatusCode != 200 {
		t.ResultWriter().Write([]byte(bodyBytes))
		execution.Status = models.Failed
		execution.Error = string(bodyBytes)
		execution.ErrorCode = res.StatusCode
		execution.CompletedAt = now
		dbactions.UpdateExecution(&execution)
		dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Failed, LastErrAt: now, LastErr: string(bodyBytes)})
		return errors.New(string(bodyBytes))
	}
	if err != nil {
		t.ResultWriter().Write([]byte(bodyBytes))
		execution.Status = models.Failed
		execution.Error = string(bodyBytes)
		execution.ErrorCode = res.StatusCode
		execution.CompletedAt = now
		dbactions.UpdateExecution(&execution)
		dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Failed, LastErrAt: now, LastErr: string(bodyBytes)})
		fmt.Printf("error here")
		return err
	}

	defer res.Body.Close()
	execution.Status = models.Success
	execution.ErrorCode = res.StatusCode
	execution.SuccessLog = string(bodyBytes)
	execution.CompletedAt = now
	dbactions.UpdateExecution(&execution)
	dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Success, CompletedAt: now, SuccessLog: string(bodyBytes)})
	return nil

}

func HandleCronCallBackTask(ctx context.Context, t *asynq.Task) error {
	var p CronPayload
	fmt.Println("started post job")
	taskId := t.ResultWriter().TaskID()
	execution := models.Execution{
		TaskID: taskId,
		Status: models.Active,
		ID:     xid.New().String(),
		RanAt:  time.Now().String(),
	}

	dbactions.CreateExecution(&execution)
	dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Active})

	fmt.Println("Started post job")
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	bodyReader := bytes.NewReader(p.Body)
	client := &http.Client{}
	req, err := http.NewRequest("POST", p.URL, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	for _, v := range p.Headers {
		headerSlice := strings.Split(v, ":")
		req.Header.Add(headerSlice[0], headerSlice[1])
	}
	res, err := client.Do(req)

	bodyBytes, _ := io.ReadAll(res.Body)
	now := time.Now().String()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.ResultWriter().Write([]byte(bodyBytes))
		execution.Status = models.Failed
		execution.Error = string(bodyBytes)
		execution.ErrorCode = res.StatusCode
		execution.CompletedAt = now
		dbactions.UpdateExecution(&execution)
		dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Failed, LastErrAt: now, LastErr: string(bodyBytes)})
		return errors.New(string(bodyBytes))
	}
	if err != nil {
		t.ResultWriter().Write([]byte(bodyBytes))
		execution.Status = models.Failed
		execution.Error = string(bodyBytes)
		execution.ErrorCode = res.StatusCode
		execution.CompletedAt = now
		dbactions.UpdateExecution(&execution)
		dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Failed, LastErrAt: now, LastErr: string(bodyBytes)})
		fmt.Printf("error here")
		return err
	}
	execution.Status = models.Success
	execution.ErrorCode = res.StatusCode
	execution.SuccessLog = string(bodyBytes)
	execution.CompletedAt = now
	dbactions.UpdateExecution(&execution)
	dbactions.UpdateTask(&models.Task{ID: taskId, Status: models.Success, CompletedAt: now, SuccessLog: string(bodyBytes)})

	log.Printf("%v res from handler", string(bodyBytes))

	return nil

}
