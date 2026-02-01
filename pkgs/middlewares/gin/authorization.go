package middleware

import (
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gin_comp"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/types"
	"github.com/gin-gonic/gin"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		v := c.Request.Context().Value(constants.ContextKeyUserInfo)
		if v == nil {
			gin_comp.ResponseError(c, base.NewUnauthorizedError("authentication required"))
			c.Abort()
			return
		}
		user, ok := v.(*types.UserAuthenticated)
		if !ok {
			gin_comp.ResponseError(c, base.NewUnauthorizedError("authentication required"))
			c.Abort()
			return
		}
		if _, ok := allowed[user.GetRole()]; !ok {
			gin_comp.ResponseError(c, base.NewForbiddenError("insufficient permissions"))
			c.Abort()
			return
		}
		c.Next()
	}
}
