package http

import (
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gin_comp"
	"github.com/gin-gonic/gin"
)

func (h *Http) HandlerUpdateBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var dto domain.DTOCreateBlog
		if err := c.ShouldBindJSON(&dto); err != nil {
			gin_comp.ResponseError(c, err)
			return
		}
		ctx := c.Request.Context()
		response, err := h.updateBlogCommand.Execute(ctx, id, &dto)
		if err != nil {
			gin_comp.ResponseError(c, err)
			return
		}
		gin_comp.ResponseSuccess(c, response)
	}
}
