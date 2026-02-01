package gorm_comp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gorm_comp/dialets"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/logger"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/opentelemetry/tracing"
)

type GormDBType int

const (
	GormDBTypeMySQL GormDBType = iota + 1
	GormDBTypePostgres
	GormDBTypeSQLite
	GormDBTypeMSSQL
	GormDBTypeNotSupported
)

type customLogger struct {
	defaultLogger logger.Logger
	slowThreshold time.Duration
	logLevel      gorm_logger.LogLevel
}

func newCustomLogger(defaultLogger logger.Logger, logLevel gorm_logger.LogLevel, slowThreshold time.Duration) *customLogger {
	return &customLogger{
		defaultLogger: defaultLogger,
		logLevel:      logLevel,
		slowThreshold: slowThreshold,
	}
}

func (l *customLogger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *customLogger) getLogger(ctx context.Context) logger.Logger {
	if ctx == nil {
		return l.defaultLogger
	}
	reqLogger := logger.FromContext(ctx)
	if reqLogger != nil {
		return reqLogger
	}
	return l.defaultLogger
}

func (l *customLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel < gorm_logger.Info {
		return
	}
	log := l.getLogger(ctx)
	args := append([]interface{}{"[GORM] " + msg}, data...)
	log.Info(args...)
}

func (l *customLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel < gorm_logger.Warn {
		return
	}
	log := l.getLogger(ctx)
	args := append([]interface{}{"[GORM] " + msg}, data...)
	log.Warn(args...)
}

func (l *customLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel < gorm_logger.Error {
		return
	}
	log := l.getLogger(ctx)
	args := append([]interface{}{"[GORM] " + msg}, data...)
	log.Error(args...)
}

func (l *customLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gorm_logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	log := l.getLogger(ctx)

	switch {
	case err != nil && l.logLevel >= gorm_logger.Error:
		log.Errorf("[GORM] SQL error [elapsed:%v] [rows:%d] sql=%s error=%v",
			elapsed, rows, sql, err)

	case elapsed > l.slowThreshold && l.slowThreshold > 0 && l.logLevel >= gorm_logger.Warn:
		log.Warnf("[GORM] slow SQL [elapsed:%v] [rows:%d] sql=%s",
			elapsed, rows, sql)

	case l.logLevel >= gorm_logger.Info:
		log.Infof("[GORM] SQL [elapsed:%v] [rows:%d] sql=%s",
			elapsed, rows, sql)
	}
}

type GormOpt struct {
	GormPrefix            string
	Dsn                   string
	DsnShadow             string
	DsnSlaves             string
	DbType                string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionIdleTime int
	EnableQueryLogging    bool
	EnableTracing         bool
}

type GormDB struct {
	logger logger.Logger
	db     *gorm.DB
	*GormOpt
}

type NewGormDBParams struct {
	fx.In
	Config *GormOpt
	Logger logger.Logger
}

func NewGormDB(p NewGormDBParams) *GormDB {
	gormDB := &GormDB{
		GormOpt: p.Config,
		logger:  p.Logger,
	}
	gormDB.validateInput()

	if err := gormDB.Activate(); err != nil {
		panic(err)
	}

	return gormDB
}

func (gdb *GormDB) validateInput() {
	errs := []string{}
	if gdb.Dsn == "" {
		errs = append(errs, "db-dsn is required")
	}
	if gdb.DbType == "" {
		errs = append(errs, "db-driver is required")
	}
	if gdb.MaxOpenConnections == 0 {
		errs = append(errs, "db-max-conn is required")
	}

	if gdb.MaxIdleConnections == 0 {
		errs = append(errs, "db-max-idle-conn is required")
	}
	if gdb.MaxConnectionIdleTime == 0 {
		errs = append(errs, "db-max-conn-idle-time is required")
	}

	if len(errs) > 0 {
		gdb.logger.Errorf("validation errors: %s", strings.Join(errs, ", "))
	}
}

func (gdb *GormDB) Activate() error {
	gdb.logger.Info("Activating GormDB")

	dbType := getDBType(gdb.DbType)
	if dbType == GormDBTypeNotSupported {
		return errors.WithStack(errors.New(fmt.Sprintf("Database type %s not supported.", gdb.DbType)))
	}

	gdb.logger.Info("Connecting to database...")

	var err error
	gdb.db, err = gdb.getDBConn(dbType)

	if err != nil {
		gdb.logger.Error("Cannot connect to database", err.Error())
		return err
	}

	return nil
}

func (gdb *GormDB) Stop() error {
	return nil
}

func (gdb *GormDB) GetDB() *gorm.DB {
	// Check if database is initialized
	if gdb.db == nil {
		gdb.logger.Fatal("Database not initialized. Make sure to call Activate() before GetDB()")
	}

	return gdb.db.Session(&gorm.Session{NewDB: true})
}

func (gdb *GormDB) GetDBWithAuditContext(ctx context.Context) *gorm.DB {
	if ctx == nil {
		ctx = context.Background()
	}

	return gdb.GetDB().WithContext(ctx)
}

func getDBType(dbType string) GormDBType {
	switch strings.ToLower(dbType) {
	case "mysql":
		return GormDBTypeMySQL
	case "postgres":
		return GormDBTypePostgres
	case "sqlite":
		return GormDBTypeSQLite
	case "mssql":
		return GormDBTypeMSSQL
	}

	return GormDBTypeNotSupported
}

func (gdb *GormDB) withInstrumentation(db *gorm.DB) {
	tracerProvider := otel.GetTracerProvider()
	if tracerProvider == nil {
		gdb.logger.Warn("Tracer provider not initialized yet, GORM OpenTelemetry plugin may not work correctly")
	} else if _, ok := tracerProvider.(*trace.TracerProvider); !ok {
		gdb.logger.Warn("Tracer provider is not a valid TracerProvider instance, GORM OpenTelemetry plugin may not work correctly")
	}

	if err := db.Use(tracing.NewPlugin(
		tracing.WithDBSystem(gdb.DbType),
		tracing.WithQueryFormatter(func(query string) string {
			return query
		}),
	)); err != nil {
		gdb.logger.Error("failed to register OpenTelemetry plugin", err)
		return
	}

	if tracerProvider != nil {
		if _, ok := tracerProvider.(*trace.TracerProvider); ok {
			gdb.logger.Info("OpenTelemetry tracing enabled for GORM", "db_system", gdb.DbType)
		} else {
			gdb.logger.Info("OpenTelemetry tracing plugin registered for GORM", "db_system", gdb.DbType)
		}
	} else {
		gdb.logger.Info("OpenTelemetry tracing plugin registered for GORM (tracer provider will be used when initialized)", "db_system", gdb.DbType)
	}
}

func (gdb *GormDB) getDBConn(t GormDBType) (dbConn *gorm.DB, err error) {
	// Parse slave DSNs if provided
	slaveDSNs := gdb.parseSlaveDSNs()

	// Determine logging level based on EnableQueryLogging flag

	var logLevel gorm_logger.LogLevel
	if gdb.EnableQueryLogging {
		logLevel = gorm_logger.Info
	} else {
		logLevel = gorm_logger.Silent
	}

	// Create custom logger with app_context.Logger
	customGormLogger := newCustomLogger(gdb.logger, logLevel, time.Second)
	gormConfig := &gorm.Config{
		Logger: customGormLogger,
	}

	// Open primary connection (source for writes)
	var db *gorm.DB
	switch t {
	case GormDBTypeMySQL:
		db, err = dialets.MySqlDB(gdb.Dsn, gormConfig)
	case GormDBTypePostgres:
		db, err = dialets.PostgresDB(gdb.Dsn, gormConfig)
	case GormDBTypeSQLite:
		db, err = dialets.SQLiteDB(gdb.Dsn, gormConfig)
	case GormDBTypeMSSQL:
		db, err = dialets.MSSqlDB(gdb.Dsn, gormConfig)
	default:
		gdb.logger.Error("unsupported database type", t)
		return nil, errors.New("unsupported database type")
	}

	if err != nil {
		gdb.logger.Error("error getting db conn", err)
		return nil, errors.WithStack(err)
	}

	gdb.logger.Infof("Connected to master database (write)")
	if len(slaveDSNs) > 0 {
		gdb.logger.Infof("Configured %d slave database(s) (read)", len(slaveDSNs))
	}

	if gdb.EnableTracing {
		gdb.withInstrumentation(db)
	}

	registerAuditHook(db)

	// Configure connection pool settings on primary connection
	// These will apply when DBResolver is not used
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(gdb.MaxOpenConnections)
		sqlDB.SetMaxIdleConns(gdb.MaxIdleConnections)
		sqlDB.SetConnMaxIdleTime(time.Duration(gdb.MaxConnectionIdleTime) * time.Second)
	}

	// If slaves are configured, set up DBResolver for read/write splitting
	if len(slaveDSNs) > 0 {
		// Create dialectors for slaves
		var slaveDialectors []gorm.Dialector
		for _, slaveDSN := range slaveDSNs {
			switch t {
			case GormDBTypeMySQL:
				slaveDialectors = append(slaveDialectors, dialets.MySqlDialector(slaveDSN))
			case GormDBTypePostgres:
				slaveDialectors = append(slaveDialectors, dialets.PostgresDialector(slaveDSN))
			case GormDBTypeSQLite:
				slaveDialectors = append(slaveDialectors, dialets.SQLiteDialector(slaveDSN))
			case GormDBTypeMSSQL:
				slaveDialectors = append(slaveDialectors, dialets.MSSqlDialector(slaveDSN))
			default:

				return nil, errors.New("unsupported database type")
			}
		}

		if len(slaveDialectors) > 0 {
			// Configure DBResolver with source (write) and replicas (read)
			// DBResolver will manage connection pools for all sources and replicas
			var masterDialector gorm.Dialector
			switch t {
			case GormDBTypeMySQL:
				masterDialector = dialets.MySqlDialector(gdb.Dsn)
			case GormDBTypePostgres:
				masterDialector = dialets.PostgresDialector(gdb.Dsn)
			case GormDBTypeSQLite:
				masterDialector = dialets.SQLiteDialector(gdb.Dsn)
			case GormDBTypeMSSQL:
				masterDialector = dialets.MSSqlDialector(gdb.Dsn)
			}

			if masterDialector != nil {
				err = db.Use(dbresolver.Register(dbresolver.Config{
					Sources:           []gorm.Dialector{masterDialector},
					Replicas:          slaveDialectors,
					Policy:            dbresolver.RandomPolicy{},
					TraceResolverMode: gdb.EnableQueryLogging,
				}).SetConnMaxIdleTime(time.Duration(gdb.MaxConnectionIdleTime) * time.Second).
					SetConnMaxLifetime(time.Duration(gdb.MaxConnectionIdleTime) * time.Second).
					SetMaxIdleConns(gdb.MaxIdleConnections).
					SetMaxOpenConns(gdb.MaxOpenConnections))
				if err != nil {
					return nil, errors.Wrap(err, "failed to register DBResolver")
				}
			}
		}
	}
	return db, nil
}

// parseSlaveDSNs parses comma-separated slave DSNs into a slice
func (gdb *GormDB) parseSlaveDSNs() []string {
	if gdb.DsnSlaves == "" {
		return nil
	}

	slaves := strings.Split(gdb.DsnSlaves, ",")
	var result []string
	for _, slave := range slaves {
		slave = strings.TrimSpace(slave)
		if slave != "" {
			result = append(result, slave)
		}
	}

	return result
}
