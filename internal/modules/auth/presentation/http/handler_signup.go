package http

import (
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/gin_comp"
	"github.com/gin-gonic/gin"
)

func (h *Http) HandlerSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			dto domain.DTOSignup
			ctx = c.Request.Context()
		)

		if err := c.ShouldBindJSON(&dto); err != nil {
			gin_comp.ResponseError(c, err)
			return
		}

		response, err := h.signupCommand.Execute(ctx, &dto)
		if err != nil {
			gin_comp.ResponseError(c, err)
			return
		}

		gin_comp.ResponseSuccessCreated(c, response)
	}
}
