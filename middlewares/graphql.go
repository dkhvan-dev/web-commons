package middlewares

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/dkhvan-dev/web-commons/constants"
	"github.com/gin-gonic/gin"
)

func GraphQLMiddleware(h *handler.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		ctx := context.WithValue(c.Request.Context(), constants.ACCEPT_LANGUAGE, c.GetHeader(constants.ACCEPT_LANGUAGE))
		c.Request = c.Request.WithContext(ctx)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
