package logger

import (
	"go.uber.org/fx/fxevent"
)

type fxEventLogger struct {
	logger Logger
}

func (l *fxEventLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Debugf("OnStart hook executing: %s", e.FunctionName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.Errorf("OnStart hook failed: %s, error: %v", e.FunctionName, e.Err)
		} else {
			l.logger.Debugf("OnStart hook executed: %s", e.FunctionName)
		}
	case *fxevent.OnStopExecuting:
		l.logger.Debugf("OnStop hook executing: %s", e.FunctionName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.Errorf("OnStop hook failed: %s, error: %v", e.FunctionName, e.Err)
		} else {
			l.logger.Debugf("OnStop hook executed: %s", e.FunctionName)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.logger.Errorf("Supplied failed: %s, error: %v", e.TypeName, e.Err)
		} else {
			l.logger.Debugf("Supplied: %s", e.TypeName)
		}
	case *fxevent.Provided:
		if e.Err != nil {
			l.logger.Errorf("Provided failed: %s, error: %v", e.ConstructorName, e.Err)
		} else {
			l.logger.Debugf("Provided: %s", e.ConstructorName)
		}
	case *fxevent.Invoking:
		l.logger.Debugf("Invoking: %s", e.FunctionName)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.logger.Errorf("Invoked failed: %s, error: %v", e.FunctionName, e.Err)
		} else {
			l.logger.Debugf("Invoked: %s", e.FunctionName)
		}
	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.Errorf("Stopped with error: %v", e.Err)
		} else {
			l.logger.Debug("Stopped")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.logger.Errorf("Started with error: %v", e.Err)
		} else {
			l.logger.Debug("Started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.Errorf("Logger initialization failed: %v", e.Err)
		} else {
			l.logger.Debugf("Logger initialized: %s", e.ConstructorName)
		}
	}
}

func ProvideFXEventLogger(logger Logger) fxevent.Logger {
	return &fxEventLogger{logger: logger}
}
