package gamelog

import (
	"fmt"
	"github.com/2pgcn/gameim/conf"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"runtime"
)

// DefaultMessageKey default message key.
var DefaultMessageKey = conf.ServerName

type GameLog interface {
	log.Logger
	LogPrint
	AppendPrefix(s string) GameLog
	NsqLog
}

type NsqLog interface {
	Output(calldepth int, s string) error
}

type LogPrint interface {
	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Warn(a ...interface{})
	Warnf(format string, a ...interface{})
	Error(a ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
}

// Helper is a logger helper.
type LogHelper struct {
	logger  log.Logger
	msgKey  string
	sprint  func(...interface{}) string
	sprintf func(format string, a ...interface{}) string
}

// NewHelper new a logger helper.
func NewHelper(logger log.Logger) *LogHelper {
	options := &LogHelper{
		msgKey:  DefaultMessageKey, // default message key
		logger:  logger,
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}
	InitGlobeLog(options)
	InitKafkaLog(options)
	defaultHelperLogger = options
	//for _, o := range opts {
	//	o(options)
	//}
	return options
}

func (h *LogHelper) AppendPrefix(s string) GameLog {
	return &LogHelper{
		msgKey:  fmt.Sprintf("%s:%s", h.msgKey, s),
		logger:  h.logger,
		sprint:  h.sprint,
		sprintf: h.sprintf,
	}
}

// Log Print log by level and keyvals.
func (h *LogHelper) Log(level log.Level, keyvals ...interface{}) error {
	_, path, file, ok := runtime.Caller(1)
	fmt.Println(path, file, ok)
	if ok {
		_ = h.logger.Log(level, "stack:", path, file)
	}
	_ = h.logger.Log(level, keyvals...)
	return nil
}

// Debug logs a message at debug level.
func (h *LogHelper) Debug(a ...interface{}) {
	_ = h.logger.Log(log.LevelDebug, h.msgKey, h.sprint(a...))
}

// Debugf logs a message at debug level.
func (h *LogHelper) Debugf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelDebug, h.msgKey, h.sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *LogHelper) Debugw(keyvals ...interface{}) {
	_ = h.logger.Log(log.LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *LogHelper) Info(a ...interface{}) {
	_ = h.logger.Log(log.LevelInfo, h.msgKey, h.sprint(a...))
}

// Infof logs a message at info level.
func (h *LogHelper) Infof(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelInfo, h.msgKey, h.sprintf(format, a...))
}

// Warn logs a message at warn level.
func (h *LogHelper) Warn(a ...interface{}) {
	_ = h.logger.Log(log.LevelWarn, h.msgKey, h.sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *LogHelper) Warnf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelWarn, h.msgKey, h.sprintf(format, a...))
}

// Error logs a message at error level.
func (h *LogHelper) Error(a ...interface{}) {
	_ = h.logger.Log(log.LevelError, h.msgKey, h.sprint(a...))
}

// Errorf logs a message at error level.
func (h *LogHelper) Errorf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelError, h.msgKey, h.sprintf(format, a...))
}

// Fatal logs a message at fatal level.
func (h *LogHelper) Fatal(a ...interface{}) {
	_ = h.logger.Log(log.LevelFatal, h.msgKey, h.sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *LogHelper) Fatalf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelFatal, h.msgKey, h.sprintf(format, a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *LogHelper) Output(calldepth int, s string) error {
	_ = h.logger.Log(log.Level(calldepth-1), fmt.Sprintf("%s-nsq", h.msgKey), h.sprintf(s))
	return nil
}
