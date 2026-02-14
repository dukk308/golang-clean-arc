package http

import (
	"net/http"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/application"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"github.com/gin-gonic/gin"
)

type GoogleOAuthHttp struct {
	googleSigninCommand *application.GoogleSigninCommand
	googleOAuthService  domain.IGoogleOAuthService
}

func NewGoogleOAuthHttp(
	googleSigninCommand *application.GoogleSigninCommand,
	googleOAuthService domain.IGoogleOAuthService,
) *GoogleOAuthHttp {
	return &GoogleOAuthHttp{
		googleSigninCommand: googleSigninCommand,
		googleOAuthService:  googleOAuthService,
	}
}

// HandlerGoogleSignin handles the Google OAuth sign-in callback
// POST /v1/auth/google/signin
func (h *GoogleOAuthHttp) HandlerGoogleSignin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto domain.DTOGoogleSignin
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := h.googleSigninCommand.Execute(c.Request.Context(), &dto)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// HandlerGoogleAuthURL returns the Google OAuth authorization URL
// GET /v1/auth/google/url
func (h *GoogleOAuthHttp) HandlerGoogleAuthURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		authURL := h.googleOAuthService.GetAuthURL(state)
		c.JSON(http.StatusOK, gin.H{"url": authURL})
	}
}
