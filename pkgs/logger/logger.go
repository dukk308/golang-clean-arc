package logger

import (
	"context"
	"os"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogType string

var (
	ZapLogType LogType = "zap"
)

type Fields map[string]interface{}

var defaultLogger Logger

// noOpLogger is a simple no-op logger used as fallback
type noOpLogger struct{}

func (n *noOpLogger) Configure(cfg func(internalLog interface{})) {}
func (n *noOpLogger) Debug(args ...interface{})                   {}
func (n *noOpLogger) Debugf(template string, args ...interface{}) {}
func (n *noOpLogger) Debugw(msg string, fields Fields)            {}
func (n *noOpLogger) LogType() LogType                            { return ZapLogType }
func (n *noOpLogger) Info(args ...interface{})                    {}
func (n *noOpLogger) Infof(template string, args ...interface{})  {}
func (n *noOpLogger) Infow(msg string, fields Fields)             {}
func (n *noOpLogger) Warn(args ...interface{})                    {}
func (n *noOpLogger) Warnf(template string, args ...interface{})  {}
func (n *noOpLogger) WarnMsg(msg string, err error)               {}
func (n *noOpLogger) Error(args ...interface{})                   {}
func (n *noOpLogger) Errorw(msg string, fields Fields)            {}
func (n *noOpLogger) Errorf(template string, args ...interface{}) {}
func (n *noOpLogger) Err(msg string, err error)                   {}
func (n *noOpLogger) Fatal(args ...interface{})                   {}
func (n *noOpLogger) Fatalf(template string, args ...interface{}) {}
func (n *noOpLogger) Printf(template string, args ...interface{}) {}
func (n *noOpLogger) WithName(name string)                        {}
func (n *noOpLogger) WithContext(ctx context.Context) Logger      { return n }
func (n *noOpLogger) WithUserInfo(ctx context.Context) Logger     { return n }

// SetDefaultLoggerInitFunc sets the function to initialize the default logger
// This allows zap_logger package to override the initialization without creating import cycle
var initDefaultLogger = func() Logger {
	logWriter := zapcore.AddSync(os.Stdout)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.NameKey = "service"
	encoderCfg.TimeKey = "time"
	encoderCfg.LevelKey = "level"
	encoderCfg.FunctionKey = "caller"
	encoderCfg.CallerKey = "line"
	encoderCfg.MessageKey = "message"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeName = zapcore.FullNameEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(zapcore.InfoLevel))

	var options []zap.Option
	options = append(options, zap.AddCaller())
	options = append(options, zap.AddCallerSkip(1))

	logger := zap.New(core, options...)
	logger = logger.Named("DEFAULT")

	return &defaultZapLoggerWrapper{
		logger:      logger,
		sugarLogger: logger.Sugar(),
	}
}

func SetDefaultLoggerInitFunc(fn func() Logger) {
	initDefaultLogger = fn
}

func init() {
	defaultLogger = initDefaultLogger()
}

// defaultZapLoggerWrapper is a simple wrapper for default zap logger
type defaultZapLoggerWrapper struct {
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
}

func (d *defaultZapLoggerWrapper) Configure(cfg func(internalLog interface{})) {
	cfg(d.logger)
}

func (d *defaultZapLoggerWrapper) Debug(args ...interface{}) {
	d.sugarLogger.Debug(args...)
}

func (d *defaultZapLoggerWrapper) Debugf(template string, args ...interface{}) {
	d.sugarLogger.Debugf(template, args...)
}

func (d *defaultZapLoggerWrapper) Debugw(msg string, fields Fields) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	d.logger.Debug(msg, zapFields...)
}

func (d *defaultZapLoggerWrapper) LogType() LogType {
	return ZapLogType
}

func (d *defaultZapLoggerWrapper) Info(args ...interface{}) {
	d.sugarLogger.Info(args...)
}

func (d *defaultZapLoggerWrapper) Infof(template string, args ...interface{}) {
	d.sugarLogger.Infof(template, args...)
}

func (d *defaultZapLoggerWrapper) Infow(msg string, fields Fields) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	d.logger.Info(msg, zapFields...)
}

func (d *defaultZapLoggerWrapper) Warn(args ...interface{}) {
	d.sugarLogger.Warn(args...)
}

func (d *defaultZapLoggerWrapper) Warnf(template string, args ...interface{}) {
	d.sugarLogger.Warnf(template, args...)
}

func (d *defaultZapLoggerWrapper) WarnMsg(msg string, err error) {
	d.logger.Warn(msg, zap.Error(err))
}

func (d *defaultZapLoggerWrapper) Error(args ...interface{}) {
	d.sugarLogger.Error(args...)
}

func (d *defaultZapLoggerWrapper) Errorw(msg string, fields Fields) {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	d.logger.Error(msg, zapFields...)
}

func (d *defaultZapLoggerWrapper) Errorf(template string, args ...interface{}) {
	d.sugarLogger.Errorf(template, args...)
}

func (d *defaultZapLoggerWrapper) Err(msg string, err error) {
	d.logger.Error(msg, zap.Error(err))
}

func (d *defaultZapLoggerWrapper) Fatal(args ...interface{}) {
	d.sugarLogger.Fatal(args...)
}

func (d *defaultZapLoggerWrapper) Fatalf(template string, args ...interface{}) {
	d.sugarLogger.Fatalf(template, args...)
}

func (d *defaultZapLoggerWrapper) Printf(template string, args ...interface{}) {
	d.sugarLogger.Infof(template, args...)
}

func (d *defaultZapLoggerWrapper) WithName(name string) {
	d.logger = d.logger.Named(name)
	d.sugarLogger = d.sugarLogger.Named(name)
}

func (d *defaultZapLoggerWrapper) WithContext(ctx context.Context) Logger {
	return d
}

func (d *defaultZapLoggerWrapper) WithUserInfo(ctx context.Context) Logger {
	return d
}

type Logger interface {
	Configure(cfg func(internalLog interface{}))
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, fields Fields)
	LogType() LogType
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, fields Fields)
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorw(msg string, fields Fields)
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
	WithName(name string)
	WithContext(ctx context.Context) Logger
	WithUserInfo(ctx context.Context) Logger
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(constants.ContextKeyRequestLogger).(Logger)
	if !ok {
		return defaultLogger
	}
	return logger
}

func ToContext(ctx context.Context, logger Logger) context.Context {
	newLogger := logger.WithContext(ctx)

	return context.WithValue(
		ctx,
		constants.ContextKeyRequestLogger,
		newLogger,
	)
}

func ToUserContext(ctx context.Context, logger Logger) context.Context {
	newLogger := logger.WithUserInfo(ctx)

	return context.WithValue(
		ctx,
		constants.ContextKeyRequestLogger,
		newLogger,
	)
}
