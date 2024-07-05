package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	gronx "github.com/adhocore/gronx"
	dbactions "github.com/dopeCape/schduler/internal/db_actions"
	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
	"github.com/dopeCape/schduler/util"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/xid"
)

func HandleNewSchdule(c *gin.Context) {
	var body struct {
		URL            string            `json:"URL"`
		Body           json.RawMessage   `json:"body"`
		Headers        map[string]string `json:"headers"`
		CronExpression string            `json:"cron"`
	}
	key := c.GetHeader("X-API-KEY")
	err := c.ShouldBindJSON(&body)
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
	if !IsCronJobExpValid(body.CronExpression) {
		c.JSON(400, gin.H{"message": "Invalide cron expression", "details": "provided cron expression is incorrent"})
		return
	}

	user, err := dbactions.GetUserFromAPIKey(key)

	uniqueId := xid.New()
	var headers []string
	for k, v := range body.Headers {
		headers = append(headers, fmt.Sprintf("%v:%v", k, v))
	}
	//handle next here
	task := models.Task{
		ID:            uniqueId.String(),
		Payload:       string(body.Body),
		Headers:       headers,
		URL:           body.URL,
		Queue:         "default",
		Retried:       0,
		Status:        models.Active,
		IsCron:        true,
		CronExpresion: body.CronExpression,
		UserID:        user.ID,
	}
	_, err = dbactions.CreateTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "details": "There was a unexpected error while enqueing the task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task info": task,
	})
	return
}

func HandleDeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing param id"})
		return
	}
	db, err := rdb.GetDb()
	dbactions.DeleteExecutionsForTask(id)
	dbactions.DeleteTask(&models.Task{ID: id})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	res := db.Delete(&models.Task{ID: id, IsCron: true})
	if res.Error != nil {
		fmt.Println(res.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shedule deleted successfully"})

}

func IsCronJobExpValid(expression string) bool {
	gron := gronx.New()
	return gron.IsValid(expression)
}

func HandleGetSchedules(c *gin.Context) {
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
	tasks, err := dbactions.GetScheduleFromApiKey(prefix, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return

	}
	c.JSON(200, gin.H{"message": tasks})

}
