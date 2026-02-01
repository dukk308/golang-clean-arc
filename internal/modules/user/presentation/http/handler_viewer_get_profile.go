package http

import (
	auth_domain "github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gin_comp"
	"github.com/gin-gonic/gin"
)

func (h *Http) HandlerViewerGetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			gin_comp.ResponseError(c, auth_domain.ErrInvalidToken)
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			gin_comp.ResponseError(c, auth_domain.ErrInvalidToken)
			return
		}

		ctx := c.Request.Context()
		response, err := h.viewerGetProfileQuery.Execute(ctx, userIDStr)
		if err != nil {
			gin_comp.ResponseError(c, err)
			return
		}

		gin_comp.ResponseSuccess(c, response)
	}
}
