package logger

import (
	"context"
	"os"

	"github.com/dukk308/golang-clean-arch-starter/pkgs/constants"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/global_config"
	log_cfg "github.com/dukk308/golang-clean-arch-starter/pkgs/logger/config"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/utils/request_id"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	level        string
	sugarLogger  *zap.SugaredLogger
	logger       *zap.Logger
	logOptions   *log_cfg.LogOptions
	globalConfig *global_config.GlobalConfig
}

type ZapLogger interface {
	Logger
	InternalLogger() *zap.Logger
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Sync() error
}

// For mapping config logger
var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

// NewDefaultZapLoggerInstances creates default zap logger instances without config
// Returns the logger and sugared logger instances that can be used to create a default logger wrapper
func NewDefaultZapLoggerInstances() (*zap.Logger, *zap.SugaredLogger) {
	logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}
	return logger, logger.Sugar()
}

// NewZapLogger create new zap logger
func NewZapLogger(
	cfg *log_cfg.LogOptions,
	globCfg *global_config.GlobalConfig,
) ZapLogger {
	zapLogger := &zapLogger{level: globCfg.LogLevel, logOptions: cfg}
	zapLogger.initLogger(globCfg)
	zapLogger.globalConfig = globCfg

	return zapLogger
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	newLogger := l.sugarLogger

	fields := make([]zap.Field, 0)

	// reqID, _ := request_id.Value(ctx)
	// if reqID != "" {
	// 	fields = append(fields, zap.String("REQID", reqID))
	// }

	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()
	if spanCtx.IsValid() {
		traceID := spanCtx.TraceID()
		if traceID.IsValid() {
			traceIDStr := traceID.String()
			fields = append(fields, zap.String("trace_id", traceIDStr))
		}

		spanID := spanCtx.SpanID()
		if spanID.IsValid() {
			spanIDStr := spanID.String()
			fields = append(fields, zap.String("span_id", spanIDStr))
		}
	} else {
		reqID, _ := request_id.Value(ctx)
		if reqID != "" {
			fields = append(fields, zap.String("trace_id", reqID))
		}
	}

	username, isHasUsername := ctx.Value(constants.ContextKeyUsername).(string)
	if isHasUsername {
		fields = append(fields, zap.String("username", username))
	}

	userID, isHasUserID := ctx.Value(constants.ContextKeyUserID).(string)
	if isHasUserID {
		fields = append(fields, zap.String("user_id", userID))
	}

	if len(fields) > 0 {
		for _, field := range fields {
			newLogger = newLogger.With(field)
		}
	}

	return &zapLogger{
		logger:       l.logger,
		sugarLogger:  newLogger,
		logOptions:   l.logOptions,
		globalConfig: l.globalConfig,
	}
}

func (l *zapLogger) WithUserInfo(ctx context.Context) Logger {
	newLogger := l.sugarLogger

	fields := make([]zap.Field, 0)

	username, isHasUsername := ctx.Value(constants.ContextKeyUsername).(string)
	if isHasUsername {
		fields = append(fields, zap.String("username", username))
	}

	userID, isHasUserID := ctx.Value(constants.ContextKeyUserID).(string)
	if isHasUserID {
		fields = append(fields, zap.String("user_id", userID))
	}

	if len(fields) > 0 {
		for _, field := range fields {
			newLogger = newLogger.With(field)
		}
	}

	return &zapLogger{
		logger:       l.logger,
		sugarLogger:  newLogger,
		logOptions:   l.logOptions,
		globalConfig: l.globalConfig,
	}
}

// WithName add logger microservice name
func (l *zapLogger) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugarLogger = l.sugarLogger.Named(name)
}

func (l *zapLogger) InternalLogger() *zap.Logger {
	return l.logger
}

func (l *zapLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// InitLogger Init logger
func (l *zapLogger) initLogger(globCfg *global_config.GlobalConfig) {
	logLevel := l.getLoggerLevel()
	logWriter := zapcore.AddSync(os.Stdout)
	isNotLocal := globCfg.Environment != "local"

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if isNotLocal {
		encoderCfg = zap.NewProductionEncoderConfig()
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
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.NameKey = "service"
		encoderCfg.TimeKey = "time"
		encoderCfg.LevelKey = "level"
		encoderCfg.FunctionKey = "caller"
		encoderCfg.CallerKey = "line"
		encoderCfg.MessageKey = "message"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		encoderCfg.ConsoleSeparator = " | "
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))

	var options []zap.Option

	if globCfg.CallerEnabled {
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddCallerSkip(1))
	}

	logger := zap.New(core, options...)

	if globCfg.EnableTracing {
		logger = otelzap.New(logger).Logger
	}

	logger = logger.Named(globCfg.ServiceName)
	l.logger = logger
	l.sugarLogger = logger.Sugar()
}

func (l *zapLogger) Configure(cfg func(internalLog interface{})) {
	cfg(l.logger)
}

func (l *zapLogger) LogType() LogType {
	return ZapLogType
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message
func (l *zapLogger) Debugf(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.Debugf(template, args...)
}

func (l *zapLogger) Debugw(msg string, fields Fields) {
	msg = buildTemplate(msg)
	zapFields := mapToZapFields(fields)
	l.logger.Debug(msg, zapFields...)
}

// Info uses fmt.Sprint to construct and log a message
func (l *zapLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Infof(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.Infof(template, args...)
}

// Infow logs a message with some additional context.
func (l *zapLogger) Infow(msg string, fields Fields) {
	msg = buildTemplate(msg)
	zapFields := mapToZapFields(fields)
	l.logger.Info(msg, zapFields...)
}

// Printf uses fmt.Sprintf to log a templated message
func (l *zapLogger) Printf(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.Infof(template, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// WarnMsg log error message with warn level.
func (l *zapLogger) WarnMsg(msg string, err error) {
	msg = buildTemplate(msg)
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Warnf(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.Warnf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorw logs a message with some additional context.
func (l *zapLogger) Errorw(msg string, fields Fields) {
	msg = buildTemplate(msg)
	zapFields := mapToZapFields(fields)
	l.logger.Error(msg, zapFields...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Errorf(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.Errorf(template, args...)
}

// Err uses error to log a message.
func (l *zapLogger) Err(msg string, err error) {
	msg = buildTemplate(msg)
	l.logger.Error(msg, zap.Error(err))
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanicf(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.DPanicf(template, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *zapLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (l *zapLogger) Panicf(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.Panicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	template = buildTemplate(template)
	l.sugarLogger.Fatalf(template, args...)
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	go func() {
		err := l.logger.Sync()
		if err != nil {
			l.logger.Error("error while syncing", zap.Error(err))
		}
	}() // nolint: errcheck
	return l.sugarLogger.Sync()
}

func mapToZapFields(data map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(data))

	for key, value := range data {
		field := convertToZapField(key, value)
		fields = append(fields, field)
	}

	return fields
}

func convertToZapField(key string, value interface{}) zap.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case int8:
		return zap.Int8(key, v)
	case int16:
		return zap.Int16(key, v)
	case int32:
		return zap.Int32(key, v)
	case int64:
		return zap.Int64(key, v)
	case uint:
		return zap.Uint(key, v)
	case uint8:
		return zap.Uint8(key, v)
	case uint16:
		return zap.Uint16(key, v)
	case uint32:
		return zap.Uint32(key, v)
	case uint64:
		return zap.Uint64(key, v)
	case bool:
		return zap.Bool(key, v)
	case float32:
		return zap.Float32(key, v)
	case float64:
		return zap.Float64(key, v)
	case error:
		return zap.NamedError(key, v)
	default:
		return zap.Any(key, v)
	}
}

func buildTemplate(template string) string {
	return template + " "
}
