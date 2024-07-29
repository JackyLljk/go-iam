package auth

import (
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"j-iam/internal/pkg/middleware"
)

// AuthzAudience 定义 jwt audience 字段
const AuthzAudience = "j-iam"

// JWTStrategy 策略实现：JWT 认证，封装了 gin-jwt
type JWTStrategy struct {
	ginjwt.GinJWTMiddleware
}

// 确保 JWTStrategy 实现了策略集（接口）
var _ middleware.AuthStrategy = &JWTStrategy{}

// NewJWTStrategy 基于 GinJWTMiddleware 创建 jwt 策略
func NewJWTStrategy(gjwt ginjwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{gjwt}
}

// AuthFunc 定义 jwt bearer 策略为 gin 认证中间件，即实现策略集的算法（实际上是 gin-jwt 实现的）
func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}
