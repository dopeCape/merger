package v1_routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dopeCape/schduler/pkg/apikey"
	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2"
)

func RegisterApiRoute(r *gin.RouterGroup, apiKeyService *apikey.ApiKeyService) {
	r.GET("/", func(c *gin.Context) {
		email := c.GetHeader("X-EMAIl")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Header X-EMAIL is missing "})
			return
		}
		user, err := apiKeyService.GenerateKey(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate key",
				"error":   "Internal server error",
			})
			return
		}

		apiKey := os.Getenv("RESEND_APIKEY")
		client := resend.NewClient(apiKey)
		params := &resend.SendEmailRequest{
			From:    "Tejes <dev@resend.dev>",
			To:      []string{email},
			Html:    fmt.Sprintf("<div>Api key <strong>%v</strong></div>", user.Key),
			Subject: "Secure api key",
		}

		_, err = client.Emails.Send(params)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate key",
				"error":   "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Api key sent to email",
		})
		return
	})

}
