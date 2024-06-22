package server

import (
	"net/http"

	v1_routes "github.com/dopeCape/schduler/internal/routes/v1"
	"github.com/dopeCape/schduler/pkg/apikey"
	"github.com/dopeCape/schduler/pkg/broker"
	"github.com/dopeCape/schduler/pkg/inspector"
	"github.com/gin-gonic/gin"
)

func GetHandler(broker *broker.Brokers, inspector *inspector.Inspector, apiKeyService *apikey.ApiKeyService) http.Handler {
	r := gin.Default()
	v1_routes.RegisterV1Routes(r, broker, inspector, apiKeyService)
	return r.Handler()

}
