package domain

func CreateViewer(dto *DTOCreateUser) (*Viewer, error) {
	viewer, err := NewUser(dto.Username, dto.Email, dto.Password)
	if err != nil {
		return nil, err
	}

	viewer.Role = RoleViewer

	return &Viewer{User: *viewer}, nil
}

func CreateAdmin(dto *DTOCreateUser) (*Admin, error) {
	admin, err := NewUser(dto.Username, dto.Email, dto.Password)
	if err != nil {
		return nil, err
	}

	admin.Role = RoleAdmin

	return &Admin{User: *admin}, nil
}

func CreateEditor(dto *DTOCreateUser) (*Editor, error) {
	editor, err := NewUser(dto.Username, dto.Email, dto.Password)
	if err != nil {
		return nil, err
	}

	editor.Role = RoleEditor

	return &Editor{User: *editor}, nil
}
