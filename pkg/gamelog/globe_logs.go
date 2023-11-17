package gamelog

import (
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"sync"
)

var defaultHelperLogger *LogHelper

var globalLogInit sync.Once

func InitGlobeLog(l *LogHelper) {
	globalLogInit.Do(func() {
		defaultHelperLogger = l
	})
}

// Log Print log by level and keyvals.
func Log(level log.Level, keyvals ...interface{}) error {
	_ = defaultHelperLogger.Log(level, keyvals...)
	return nil
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	defaultHelperLogger.Debug(a)
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelDebug, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelInfo, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelInfo, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelWarn, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelWarn, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelError, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelError, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelFatal, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	_ = defaultHelperLogger.Log(log.LevelFatal, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
	os.Exit(1)
}
