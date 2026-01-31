package dialets

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SQLiteDB Get SQLite DB connection
// dsn string
// Ex: /tmp/gorm.db
func SQLiteDB(dsn string, config *gorm.Config) (db *gorm.DB, err error) {
	if config == nil {
		config = &gorm.Config{}
	}
	return gorm.Open(sqlite.Open(dsn), config)
}

// SQLiteDialector returns a SQLite Dialector for DBResolver
// dsn string
// Ex: /tmp/gorm.db
func SQLiteDialector(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
