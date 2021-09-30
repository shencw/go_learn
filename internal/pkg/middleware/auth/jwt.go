package auth

import (
	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go_learn/internal/pkg/middleware"
)

const AuthzAudience = "127.0.0.1"

type JWTStrategy struct {
	ginJwt.GinJWTMiddleware
}

var _ middleware.AuthStrategy = &JWTStrategy{}

func NewJWTStrategy(gJwt ginJwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{gJwt}
}

func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}
