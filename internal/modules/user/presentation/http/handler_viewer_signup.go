package http

import (
	"github.com/dukk308/golang-clean-arc/internal/modules/user/domain"
	"github.com/dukk308/golang-clean-arc/pkgs/components/gin_comp"
	"github.com/gin-gonic/gin"
)

func (h *Http) HandlerViewerSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			dto domain.DTOCreateUser
			ctx = c.Request.Context()
		)

		if err := c.ShouldBindJSON(&dto); err != nil {
			gin_comp.ResponseError(c, err)
			return
		}

		response, err := h.createNewViewerCommand.Execute(ctx, &dto)
		if err != nil {
			gin_comp.ResponseError(c, err)
			return
		}

		gin_comp.ResponseSuccessCreated(c, response)
	}
}
