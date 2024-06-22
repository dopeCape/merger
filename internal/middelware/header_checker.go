package middelware

import (
	"net/http"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
	"github.com/gin-gonic/gin"
)

func KeyChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Api key missing"})

		} else {
			db, err := rdb.GetDb()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "There was a problem procecssing your requrest"})
			}

			var key models.ApiKey
			res := db.Where("key = ?", apiKey).First(&key)
			if res.Error != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid api key"})
			}

			c.Next()
		}
	}
}
