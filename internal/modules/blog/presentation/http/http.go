package http

import (
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/application"
	"github.com/gin-gonic/gin"
)

type Http struct {
	createBlogCommand *application.CreateBlogCommand
	getBlogQuery      *application.GetBlogQuery
	listBlogsQuery    *application.ListBlogsQuery
	updateBlogCommand *application.UpdateBlogCommand
	deleteBlogCommand *application.DeleteBlogCommand
}

func NewHttp(
	createBlogCommand *application.CreateBlogCommand,
	getBlogQuery *application.GetBlogQuery,
	listBlogsQuery *application.ListBlogsQuery,
	updateBlogCommand *application.UpdateBlogCommand,
	deleteBlogCommand *application.DeleteBlogCommand,
) *Http {
	return &Http{
		createBlogCommand: createBlogCommand,
		getBlogQuery:      getBlogQuery,
		listBlogsQuery:    listBlogsQuery,
		updateBlogCommand: updateBlogCommand,
		deleteBlogCommand: deleteBlogCommand,
	}
}

func (h *Http) RegisterRoutes(router *gin.RouterGroup) {
	public := router.Group("/public/v1/blogs")
	{
		public.GET("", h.HandlerListBlogs())
		public.GET("/:slug", h.HandlerGetBlogBySlug())
	}
	admin := router.Group("/admin/v1/blogs")
	{
		admin.POST("", h.HandlerCreateBlog())
		admin.GET("", h.HandlerListBlogs())
		admin.GET("/:id", h.HandlerGetBlogByID())
		admin.PUT("/:id", h.HandlerUpdateBlog())
		admin.DELETE("/:id", h.HandlerDeleteBlog())
	}
}
