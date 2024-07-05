package middelware

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func keyFunc(c *gin.Context, header string) string {
	hashedKey := sha256.Sum256([]byte(strings.Split(c.GetHeader(header), ".")[1]))
	return hex.EncodeToString(hashedKey[:])
}
func keyFuncEmail(c *gin.Context) string {
	return c.GetHeader("X-EMAIl")
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.JSON(429, "Too many requests. Try again in "+time.Until(info.ResetTime.Truncate(time.Second)).String())
}

func GetRateLimiter(key string, limit uint, reset int) gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		Rate:  (time.Minute * time.Duration(reset)),
		Limit: limit,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      func(c *gin.Context) string { return keyFunc(c, key) },
	})
	return mw
}
func GetEmailRateLimiter(limit uint, reset int) gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		Rate:  (time.Minute * time.Duration(reset)),
		Limit: limit,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      func(c *gin.Context) string { return keyFuncEmail(c) },
	})
	return mw
}
