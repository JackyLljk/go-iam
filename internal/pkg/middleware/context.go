package middleware

import (
	"github.com/gin-gonic/gin"

	"j-iam/pkg/log"
)

// UsernameKey 在 gin.Context 中定义 key，该 key 表示密钥的所有者
const UsernameKey = "username"

// Context is a middleware that injects common prefix fields to gin.Context.
// 中间件，向 gin.Context 注入 RequestIDKey 和 UsernameKey
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.KeyRequestID, c.GetString(XRequestIDKey))
		c.Set(log.KeyUsername, c.GetString(UsernameKey))
		c.Next()
	}
}
