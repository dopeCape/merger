package v1_routes

import (
	"github.com/dopeCape/schduler/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterScheduleRouter(gr *gin.RouterGroup) {
	gr.POST("/schedule", func(c *gin.Context) {
		handler.HandleNewSchdule(c)
	})
}
