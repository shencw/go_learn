package middleware

import "github.com/gin-gonic/gin"

// AuthStrategy 定义了一组用于进行资源认证的策略
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator 用于在不同策略之间进行切换
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy 用于设置另一个鉴权策略
func (o *AuthOperator) SetStrategy(strategy AuthStrategy) {
	o.strategy = strategy
}

// AuthFunc 执行鉴权资源
func (o *AuthOperator) AuthFunc() gin.HandlerFunc {
	return o.strategy.AuthFunc()
}