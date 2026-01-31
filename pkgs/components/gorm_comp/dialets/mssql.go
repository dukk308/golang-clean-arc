package dialets

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// MSSqlDB Get MS SQL DB connection
// dsn string
// Ex: sqlserver://username:password@localhost:1433?database=dbname
func MSSqlDB(dsn string, config *gorm.Config) (db *gorm.DB, err error) {
	if config == nil {
		config = &gorm.Config{}
	}
	return gorm.Open(sqlserver.Open(dsn), config)
}

// MSSqlDialector returns a MS SQL Dialector for DBResolver
// dsn string
// Ex: sqlserver://username:password@localhost:1433?database=dbname
func MSSqlDialector(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}
