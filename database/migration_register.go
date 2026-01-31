package database

import (
	user_persistence "github.com/dukk308/golang-clean-arch-starter/internal/modules/user/infrastructure/persistence"
)

var Models = []interface{}{
	user_persistence.SQLUser{},
}
