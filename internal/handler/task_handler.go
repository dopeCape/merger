package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dopeCape/schduler/pkg/broker"
	"github.com/dopeCape/schduler/pkg/inspector"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
)

func HandleEnqueue(c *gin.Context, broker *broker.Brokers) {

	var body struct {
		URL     string            `json:"URL"`
		Body    json.RawMessage   `json:"body"`
		Delay   time.Duration     `json:"delay"`
		Time    string            `json:"time"`
		Headers map[string]string `json:"headers"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"message": "Bad request", "details": "There was a error procecssing request body"})
		return
	}
	if body.URL == "" {
		c.JSON(400, gin.H{"message": "URL is missing", "details": "proptery URL is missing from request body "})
		return
	}
	if len(body.Body) == 0 {
		c.JSON(400, gin.H{"message": "body is missing", "details": "proptery body is missing from request body"})
		return
	}

	var taskinfo *asynq.TaskInfo
	if body.Delay != 0 {
		taskinfo, err = broker.EnqueueWithDelay(body.URL, body.Body, body.Headers, body.Delay)
	} else if body.Time != "" {
		taskinfo, err = broker.EnqueueAt(body.URL, body.Body, body.Headers, body.Time)
	} else {
		taskinfo, err = broker.Enqueue(body.URL, body.Body, body.Headers)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "details": "There was a unexpected error while enqueing the task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"task info": taskinfo,
	})
	return
}

func HandleDequque(c *gin.Context, inspector *inspector.Inspector) {
	queue := c.Param("queue")
	fmt.Println(queue)
	if queue == "" {
		c.JSON(400, gin.H{"message": "Bad request", "details": "paramerter quque not found "})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"message": "Bad request", "details": "paramerter quque not found "})
		return
	}
	err := inspector.Dequque(queue, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task dequed succesfully"})

}
