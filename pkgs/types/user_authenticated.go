package types

type UserAuthenticated struct {
	ID    string
	Email string
	Role  string
}

func (u *UserAuthenticated) GetID() string {
	return u.ID
}

func (u *UserAuthenticated) GetEmail() string {
	return u.Email
}

func (u *UserAuthenticated) GetRole() string {
	return u.Role
}
