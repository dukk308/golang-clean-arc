package domain

import "time"

type DTOProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func NewDTOProfileResponse(viewer *Viewer) *DTOProfileResponse {
	return &DTOProfileResponse{
		ID:        viewer.ID.String(),
		Username:  viewer.Username,
		Email:     viewer.Email.Value,
		Role:      viewer.Role.String(),
		CreatedAt: viewer.CreatedAt,
		UpdatedAt: viewer.UpdatedAt,
	}
}
