package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	dbactions "github.com/dopeCape/schduler/internal/db_actions"
	"github.com/dopeCape/schduler/internal/models"
	"github.com/dopeCape/schduler/pkg/broker"
	"github.com/dopeCape/schduler/pkg/inspector"
	"github.com/dopeCape/schduler/util"
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
	key := c.GetHeader("X-API-KEY")
	if err != nil {
		c.JSON(400, gin.H{"message": "Bad request", "details": "There was a error procecssing request body"})
		return
	}
	maxSize := os.Getenv("MAX_PAYLOAD")
	payloadSizeInCheck, err := util.PayloadSizeChecker(string(body.Body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "details": "There was a unexpected error while enqueing the task"})
		return
	}
	if !payloadSizeInCheck {
		c.JSON(400, gin.H{"message": "Bad request", "details": fmt.Sprintf("body size is capped to %v kb ", maxSize)})
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
	user, err := dbactions.GetUserFromAPIKey(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "details": "There was a unexpected error while enqueing the task"})
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
	var headers []string
	for k, v := range body.Headers {
		headers = append(headers, fmt.Sprintf("%v:%v", k, v))
	}

	task := models.Task{
		ID:      taskinfo.ID,
		Payload: string(body.Body),
		Headers: headers,
		URL:     body.URL,
		Queue:   taskinfo.Queue,
		Retried: 0,
		Next:    taskinfo.NextProcessAt.String(),
		Status:  models.Active,
		IsCron:  false,
		UserID:  user.ID,
	}
	dbactions.CreateTask(&task)
	c.JSON(http.StatusOK, gin.H{
		"task info": task,
	})
	return
}

func HandleDequque(c *gin.Context, inspector *inspector.Inspector) {
	queue := c.Param("queue")
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task is  in active or completed state"})
		return
	}

	err = dbactions.DeleteExecutionsForTask(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server error"})
		return
	}

	err = dbactions.DeleteTask(&models.Task{ID: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task dequed succesfully"})

}

func GetTasks(c *gin.Context) {
	offsetStr := c.Query("offset")
	var offset int = 0
	var limit int = 20
	prefix := strings.Split(c.GetHeader("X-API-KEY"), ".")[0]
	var err error
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
	}
	if err != nil {
		c.JSON(400, gin.H{"message": "Bad request", "details": "failed to read query param offset"})
		return
	}
	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
	}
	if err != nil {
		c.JSON(400, gin.H{"message": "Bad request", "details": "failed to read query param offset"})
		return
	}
	tasks, err := dbactions.GetTaskForAPIKey(prefix, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": tasks})
}
func RunNow(c *gin.Context, inspector *inspector.Inspector) {
	queue := "default"
	taskID := c.Query("taskID")
	if taskID == "" {
		c.JSON(400, gin.H{"message": "Bad request", "details": "Missing query param taskID"})
		return
	}
	err := inspector.RunNow(taskID, queue)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "Task run succesfully"})
	return
}
