package http

import (
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/gin_comp"
	"github.com/gin-gonic/gin"
)

func (h *Http) HandlerRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			dto domain.DTORefreshToken
			ctx = c.Request.Context()
		)

		if err := c.ShouldBindJSON(&dto); err != nil {
			gin_comp.ResponseError(c, err)
			return
		}

		response, err := h.refreshTokenCommand.Execute(ctx, &dto)
		if err != nil {
			gin_comp.ResponseError(c, err)
			return
		}

		gin_comp.ResponseSuccess(c, response)
	}
}
