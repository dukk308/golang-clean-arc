package gorm_comp

import (
	"flag"
)

var (
	dbPrefixVal string
	dbDsnVal    string
	dbDsnShadowVal string
	dbDsnSlavesVal string
	dbDriverVal string
	dbMaxConnVal int
	dbMaxIdleConnVal int
	dbMaxConnIdleTimeVal int
	dbEnableQueryLogVal bool
	dbEnableTracingVal bool
)

var (
	DbPrefix = &dbPrefixVal
	DbDsn    = &dbDsnVal
	DbDsnShadow = &dbDsnShadowVal
	DbDsnSlaves = &dbDsnSlavesVal
	DbDriver = &dbDriverVal
	DbMaxConn = &dbMaxConnVal
	DbMaxIdleConn = &dbMaxIdleConnVal
	DbMaxConnIdleTime = &dbMaxConnIdleTimeVal
	DbEnableQueryLog = &dbEnableQueryLogVal
	DbEnableTracing = &dbEnableTracingVal
)

func init() {
	if flag.Lookup("db-prefix") == nil {
		flag.StringVar(&dbPrefixVal, "db-prefix", "", "Database prefix")
	}
	if flag.Lookup("db-dsn") == nil {
		flag.StringVar(&dbDsnVal, "db-dsn", "", "Database DSN")
	}
	if flag.Lookup("db-dsn-shadow") == nil {
		flag.StringVar(&dbDsnShadowVal, "db-dsn-shadow", "", "Database shadow DSN")
	}
	if flag.Lookup("db-dsn-slaves") == nil {
		flag.StringVar(&dbDsnSlavesVal, "db-dsn-slaves", "", "Database slaves DSN")
	}
	if flag.Lookup("db-driver") == nil {
		flag.StringVar(&dbDriverVal, "db-driver", "", "Database driver")
	}
	if flag.Lookup("db-max-conn") == nil {
		flag.IntVar(&dbMaxConnVal, "db-max-conn", 0, "Maximum database connections")
	}
	if flag.Lookup("db-max-idle-conn") == nil {
		flag.IntVar(&dbMaxIdleConnVal, "db-max-idle-conn", 0, "Maximum idle database connections")
	}
	if flag.Lookup("db-max-conn-idle-time") == nil {
		flag.IntVar(&dbMaxConnIdleTimeVal, "db-max-conn-idle-time", 0, "Maximum connection idle time in seconds")
	}
	if flag.Lookup("db-enable-query-log") == nil {
		flag.BoolVar(&dbEnableQueryLogVal, "db-enable-query-log", false, "Enable database query logging")
	}
	if flag.Lookup("db-enable-tracing") == nil {
		flag.BoolVar(&dbEnableTracingVal, "db-enable-tracing", false, "Enable database tracing")
	}
}

func LoadDatabaseConfigs() *GormOpt {
	result := &GormOpt{
		GormPrefix:            *DbPrefix,
		Dsn:                   *DbDsn,
		DsnShadow:             *DbDsnShadow,
		DsnSlaves:             *DbDsnSlaves,
		DbType:                *DbDriver,
		MaxOpenConnections:    *DbMaxConn,
		MaxIdleConnections:    *DbMaxIdleConn,
		MaxConnectionIdleTime: *DbMaxConnIdleTime,
		EnableQueryLogging:    *DbEnableQueryLog,
		EnableTracing:         *DbEnableTracing,
	}

	return result
}
