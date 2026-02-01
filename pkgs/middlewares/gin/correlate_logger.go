package middleware

import (
	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"
	"github.com/gin-gonic/gin"
)

func CorrelateLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.ToContext(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
