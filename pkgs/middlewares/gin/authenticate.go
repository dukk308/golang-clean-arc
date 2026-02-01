package middleware

import (
	"context"
	"strings"

	auth_domain "github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gin_comp"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/types"
	"github.com/gin-gonic/gin"
)

func Authenticate(tokenService auth_domain.ITokenService) gin.HandlerFunc {
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

		ctx := context.WithValue(
			c.Request.Context(),
			constants.ContextKeyUserInfo,
			&types.UserAuthenticated{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  claims.Role,
			},
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
