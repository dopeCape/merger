package suscriber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"
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
	var p Payload
	fmt.Println("Started post job")

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	fmt.Println(p.Headers)
	bodyReader := bytes.NewReader(p.Body)
	client := &http.Client{}
	req, err := http.NewRequest("POST", p.URL, bodyReader)
	fmt.Println(req.Body, "body")
	req.Header.Add("Content-Type", "application/json")
	for k, v := range p.Headers {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	log.Printf("%v res from handler", string(bodyBytes))

	return nil

}

func HandleCronCallBackTask(ctx context.Context, t *asynq.Task) error {
	var p CronPayload
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
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	log.Printf("%v res from handler", string(bodyBytes))

	return nil

}
