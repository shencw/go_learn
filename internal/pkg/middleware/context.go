package middleware

import (
	"github.com/gin-gonic/gin"
	"go_learn/pkg/log"
)

const UsernameKey = "username"

func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyUsername, c.GetString(UsernameKey))
		c.Next()
	}
}
