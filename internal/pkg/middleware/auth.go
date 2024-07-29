package middleware

import (
	"github.com/gin-gonic/gin"
)

// 基于策略模式，灵活选择 basic 和 jwt 认证

// AuthStrategy 策略类：定义策略需要实现的方法 AuthFunc() gin.HandlerFunc
type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

// AuthOperator 抽象出策略执行者，切换不同的策略实现
type AuthOperator struct {
	strategy AuthStrategy
}

// SetStrategy 切换策略
func (operator *AuthOperator) SetStrategy(strategy AuthStrategy) {
	operator.strategy = strategy
}

// AuthFunc 用策略中的方法
func (operator *AuthOperator) AuthFunc() gin.HandlerFunc {
	return operator.strategy.AuthFunc()
}
