package middleware

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	"github.com/gin-gonic/gin"
)

func CorrelateLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), constants.ContextKeyRequestLogger, log)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
