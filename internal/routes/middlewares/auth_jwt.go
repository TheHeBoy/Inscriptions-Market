// Package middlewares Gin 中间件
package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
