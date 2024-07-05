package handler

import (
	dbactions "github.com/dopeCape/schduler/internal/db_actions"
	"github.com/gin-gonic/gin"
)

func HandleGetExecutions(c *gin.Context) {
	taskID := c.Query("taskID")
	if taskID == "" {
		c.JSON(400, gin.H{"message": "Bad request", "details": "missing query parameter taskID"})
		return
	}
	executions, err := dbactions.GetExecutionsForTask(taskID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"data": executions})
	return
}
