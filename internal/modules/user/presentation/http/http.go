package http

import (
	auth_domain "github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user/application"
	"github.com/gin-gonic/gin"
)

type Http struct {
	getProfileQuery *application.GetProfileQuery
	tokenService    auth_domain.ITokenService
}

func NewHttp(
	getProfileQuery *application.GetProfileQuery,
	tokenService auth_domain.ITokenService,
) *Http {
	return &Http{
		getProfileQuery: getProfileQuery,
		tokenService:    tokenService,
	}
}

func (h *Http) RegisterRoutes(router *gin.RouterGroup) {
	accountGroup := router.Group("/v1/account")
	accountGroup.Use(AuthMiddleware(h.tokenService))
	{
		accountGroup.GET("/profile", h.HandlerGetProfile())
	}
}
