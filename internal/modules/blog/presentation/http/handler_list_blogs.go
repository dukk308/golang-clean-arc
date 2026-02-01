package http

import (
	"strconv"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gin_comp"
	"github.com/gin-gonic/gin"
)

func (h *Http) HandlerListBlogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		ctx := c.Request.Context()
		response, err := h.listBlogsQuery.Execute(ctx, page, limit)
		if err != nil {
			gin_comp.ResponseError(c, err)
			return
		}
		gin_comp.ResponseSuccess(c, response)
	}
}
