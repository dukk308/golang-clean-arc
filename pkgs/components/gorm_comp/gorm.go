package gorm_comp

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormComponent struct {
	db *gorm.DB
}

func NewGormComponent() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/test"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
