package handler

import (
	"fmt"
	"net/http"

	gronx "github.com/adhocore/gronx"
	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
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
	if !IsCronJobExpValid(body.CronExpression) {
		c.JSON(400, gin.H{"message": "Invalide cron expression", "details": "provided cron expression is incorrent"})
		return
	}

	uniqueId := xid.New()
	var headers []string
	for k, v := range body.Headers {
		headers = append(headers, fmt.Sprintf("%v:%v", k, v))
	}
	task := &models.Task{ID: uniqueId.String(), URL: body.URL, Payload: string(body.Body), IsCron: true, CronExpresion: body.CronExpression, Headers: headers}
	db, err := rdb.GetDb()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "details": "There was a unexpected error while scheduling the task"})
		return
	}
	res := db.Create(&task)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "details": "There was a unexpected error while scheduling the task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task info": task,
	})
	return
}

func IsCronJobExpValid(expression string) bool {
	gron := gronx.New()
	return gron.IsValid(expression)
}
