package auth

import (
	"github.com/gin-gonic/gin"
	"go_learn/internal/pkg/code"
	"go_learn/internal/pkg/middleware"
	"go_learn/pkg/core"
	"go_learn/pkg/errors"
	"strings"
)

type AutoStrategy struct {
	basic BasicStrategy
	jwt   JWTStrategy
}

var _ middleware.AuthStrategy = &AutoStrategy{}

// NewAutoStrategy create auto strategy with basic strategy and jwt strategy.
func NewAutoStrategy(basic BasicStrategy, jwt JWTStrategy) AutoStrategy {
	return AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}

func (a AutoStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		operator := middleware.AuthOperator{}
		authHeader := strings.SplitN(c.Request.Header.Get(Authorization), " ", AuthHeaderCount)
		if len(authHeader) != AuthHeaderCount {
			core.WriteResponse(
				c, errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."), nil,
			)
			c.Abort()
			return
		}

		switch authHeader[0] {
		case Basic:
			operator.SetStrategy(a.basic)
		case Bearer:
			operator.SetStrategy(a.jwt)
		default:
			core.WriteResponse(
				c, errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header."), nil,
			)
			c.Abort()
			return
		}

		operator.AuthFunc()(c)

		c.Next()
	}
}
