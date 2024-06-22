package middelware

import (
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func keyFunc(c *gin.Context) string {
	return c.GetHeader("X-API-KEY")
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.JSON(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func GetRateLimiter() gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		Rate:  time.Minute,
		Limit: 5,
	})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
	return mw
}
