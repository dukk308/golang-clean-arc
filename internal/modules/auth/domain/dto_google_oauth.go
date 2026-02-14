package domain

// DTOGoogleSignin represents the Google OAuth sign-in request
type DTOGoogleSignin struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"`
}

// DTOGoogleUser represents the user info from Google
type DTOGoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}
