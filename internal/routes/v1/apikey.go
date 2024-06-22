package v1_routes

import (
	"net/http"

	"github.com/dopeCape/schduler/pkg/apikey"
	"github.com/gin-gonic/gin"
)

func RegisterApiRoute(r *gin.Engine, apiKeyService *apikey.ApiKeyService) {
	r.GET("/apikey/generate", func(c *gin.Context) {
		key, err := apiKeyService.GenerateKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate key",
				"error":   "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"apiKey": key,
		})
		return
	})

}
