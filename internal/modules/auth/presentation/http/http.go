package http

import (
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/application"
	"github.com/gin-gonic/gin"
)

type Http struct {
	signupCommand       *application.SignupCommand
	signinCommand       *application.SigninCommand
	signoutCommand      *application.SignoutCommand
	refreshTokenCommand *application.RefreshTokenCommand
	googleOAuthHttp     *GoogleOAuthHttp
}

func NewHttp(
	signupCommand *application.SignupCommand,
	signinCommand *application.SigninCommand,
	signoutCommand *application.SignoutCommand,
	refreshTokenCommand *application.RefreshTokenCommand,
	googleOAuthHttp *GoogleOAuthHttp,
) *Http {
	return &Http{
		signupCommand:       signupCommand,
		signinCommand:       signinCommand,
		signoutCommand:      signoutCommand,
		refreshTokenCommand: refreshTokenCommand,
		googleOAuthHttp:     googleOAuthHttp,
	}
}

func (h *Http) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/v1/auth/signup", h.HandlerSignup())
	router.POST("/v1/auth/signin", h.HandlerSignin())
	router.POST("/v1/auth/signout", h.HandlerSignout())
	router.POST("/v1/auth/refresh", h.HandlerRefreshToken())
	router.GET("/v1/auth/google/url", h.googleOAuthHttp.HandlerGoogleAuthURL())
	router.POST("/v1/auth/google/signin", h.googleOAuthHttp.HandlerGoogleSignin())
}
