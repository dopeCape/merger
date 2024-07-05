package middelware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
	"github.com/gin-gonic/gin"
)

func KeyChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		fmt.Println(apiKey)
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Api key missing"})
			return
		} else {
			keySlice := strings.Split(apiKey, ".")
			hashedKey := sha256.Sum256([]byte(keySlice[1]))
			prefix := keySlice[0]
			key := hex.EncodeToString(hashedKey[:])
			db, err := rdb.GetDb()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "There was a problem procecssing your requrest"})
				return
			}
			var user models.User
			res := db.Where("key = ? AND prefix = ?", key, prefix).First(&user)
			if res.Error != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid api key"})
				return
			}
			c.Next()
		}

	}
}
