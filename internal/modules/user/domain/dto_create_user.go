package domain

type DTOCreateUser struct {
	Username string `json:"username,omitempty"`	
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Provider AuthProvider `json:"provider"`
	ProviderID string `json:"provider_id,omitempty"`
}
