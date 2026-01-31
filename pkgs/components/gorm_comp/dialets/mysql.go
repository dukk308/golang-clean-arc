package dialets

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySqlDB Get MySQL DB connection
// dsn string
// Ex: user:password@/db_name?charset=utf8&parseTime=True&loc=Local
func MySqlDB(dsn string, config *gorm.Config) (db *gorm.DB, err error) {
	if config == nil {
		config = &gorm.Config{}
	}
	return gorm.Open(mysql.Open(dsn), config)
}

// MySqlDialector returns a MySQL Dialector for DBResolver
// dsn string
// Ex: user:password@/db_name?charset=utf8&parseTime=True&loc=Local
func MySqlDialector(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}
