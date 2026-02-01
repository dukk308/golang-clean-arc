package http

import (
	auth_domain "github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/user/application"
	"github.com/gin-gonic/gin"
)

type Http struct {
	viewerGetProfileQuery *application.ViewerGetProfileQuery
	tokenService          auth_domain.ITokenService
}

func NewHttp(
	viewerGetProfileQuery *application.ViewerGetProfileQuery,
	tokenService auth_domain.ITokenService,
) *Http {
	return &Http{
		viewerGetProfileQuery: viewerGetProfileQuery,
		tokenService:          tokenService,
	}
}

func (h *Http) RegisterRoutes(router *gin.RouterGroup) {
	accountGroup := router.Group("/v1/account")
	accountGroup.Use(AuthMiddleware(h.tokenService))
	{
		accountGroup.GET("/profile", h.HandlerViewerGetProfile())
	}
}
