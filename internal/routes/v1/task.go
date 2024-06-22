package v1_routes

import (
	"github.com/dopeCape/schduler/internal/handler"
	"github.com/dopeCape/schduler/pkg/broker"
	"github.com/dopeCape/schduler/pkg/inspector"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRouter(gr *gin.RouterGroup, broker *broker.Brokers, inspector *inspector.Inspector) {
	gr.POST("/task/enqueue", func(c *gin.Context) {
		handler.HandleEnqueue(c, broker)
	})
	gr.DELETE("/task/:queue/:id", func(c *gin.Context) {
		handler.HandleDequque(c, inspector)
	})

}
