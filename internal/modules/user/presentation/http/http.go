package http

import (
	"github.com/dukk308/golang-clean-arc/internal/modules/user/application"
	"github.com/gin-gonic/gin"
)

type Http struct {
	createNewViewerCommand *application.CreateNewViewerCommand
}

func NewHttp(
	createNewViewerCommand *application.CreateNewViewerCommand,
) *Http {
	return &Http{
		createNewViewerCommand: createNewViewerCommand,
	}
}

func (h *Http) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/v1/auth/signup", h.HandlerViewerSignup())
	router.POST("/v1/auth/signin", nil)
	router.POST("/v1/auth/signout", nil)
	router.POST("/v1/auth/refresh", nil)
}
