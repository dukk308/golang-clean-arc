package dialets

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDB Get Postgres DB connection
// dns string
// Ex: host=myhost port=myport user=gorm dbname=gorm password=mypassword
func PostgresDB(dsn string, config *gorm.Config) (db *gorm.DB, err error) {
	if config == nil {
		config = &gorm.Config{}
	}
	return gorm.Open(postgres.Open(dsn), config)
}

// PostgresDialector returns a Postgres Dialector for DBResolver
// dsn string
// Ex: host=myhost port=myport user=gorm dbname=gorm password=mypassword
func PostgresDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
