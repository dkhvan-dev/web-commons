package middlewares

import (
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config.Logger.Info("HTTP request",
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.Next()
	}
}
