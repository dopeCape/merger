package v1_routes

import (
	"github.com/dopeCape/schduler/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterExecutionRouter(gr *gin.RouterGroup) {
	gr.GET("/executions", handler.HandleGetExecutions)

}
