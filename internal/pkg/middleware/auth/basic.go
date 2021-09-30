package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go_learn/internal/pkg/code"
	"go_learn/internal/pkg/middleware"
	"go_learn/pkg/core"
	"go_learn/pkg/errors"
	"strings"
)

type BasicStrategy struct {
	compare func(username string, password string) bool
}

// NewBasicStrategy create basic strategy with compare function.
func NewBasicStrategy(compare func(username string, password string) bool) BasicStrategy {
	return BasicStrategy{
		compare: compare,
	}
}

func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get(Authorization), " ", AuthHeaderCount)

		if len(auth) != AuthHeaderCount || auth[0] != Basic {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."), nil)
			c.Abort()
			return
		}

		// 解密token
		payload, _ := base64.StdEncoding.DecodeString(auth[1])

		// pair[0]账号 pair[1]密码
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."), nil)
			c.Abort()
			return
		}

		c.Set(middleware.UsernameKey, pair[0])

		c.Next()
	}
}
