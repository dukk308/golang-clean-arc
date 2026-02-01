package database

import (
	blog_persistence "github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/infrastructure/persistence"
	note_persistence "github.com/dukk308/beetool.dev-go-starter/internal/modules/note/infrastructure/persistence"
	user_persistence "github.com/dukk308/beetool.dev-go-starter/internal/modules/user/infrastructure/persistence"
)

var Models = []interface{}{
	user_persistence.SQLUser{},
	note_persistence.SQLNote{},
	blog_persistence.SQLBlog{},
}
