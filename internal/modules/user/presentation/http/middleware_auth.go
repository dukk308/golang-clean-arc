package http

import (
	"strings"

	auth_domain "github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/gin_comp"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService auth_domain.ITokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			gin_comp.ResponseError(c, auth_domain.ErrInvalidToken)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			gin_comp.ResponseError(c, auth_domain.ErrInvalidToken)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := tokenService.ValidateToken(token)
		if err != nil {
			gin_comp.ResponseError(c, err)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}
