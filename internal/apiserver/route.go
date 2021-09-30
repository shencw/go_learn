package apiserver

import (
	"github.com/gin-gonic/gin"
	"go_learn/internal/pkg/middleware"
	"go_learn/internal/pkg/middleware/auth"
	"go_learn/pkg/core"
	"net/http"
)

func Route() http.Handler {
	router := gin.Default()

	// Middlewares.
	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	router.POST("/login", jwtStrategy.LoginHandler)
	router.POST("/logout", jwtStrategy.LogoutHandler)
	router.POST("/refresh", jwtStrategy.RefreshHandler)

	// 路由分组、中间件、认证
	v1 := router.Group("/v1")
	{
		auto := newAutoAuth()
		v1.Use(auto.AuthFunc()).GET("/login", func(c *gin.Context) {
			core.WriteResponse(c, nil, "Hi:"+c.Value(middleware.UsernameKey).(string))
		})
	}

	return router
}
