package v1_routes

import (
	"github.com/dopeCape/schduler/internal/middelware"
	"github.com/dopeCape/schduler/pkg/apikey"
	"github.com/dopeCape/schduler/pkg/broker"
	"github.com/dopeCape/schduler/pkg/inspector"
	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(r *gin.Engine, broker *broker.Brokers, inspector *inspector.Inspector, apiKeyService *apikey.ApiKeyService) {
	group := r.Group("/api/v1")
	group.Use(middelware.KeyChecker(), middelware.GetRateLimiter())
	RegisterTaskRouter(group, broker, inspector)
	RegisterApiRoute(r, apiKeyService)
	RegisterScheduleRouter(group)
}
