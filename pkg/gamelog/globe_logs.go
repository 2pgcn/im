package gamelog

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
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

func GetGlobalog() *LogHelper {
	return defaultHelperLogger
}

// todo
func GetGlobalNsqLog() *LogHelper {
	l := GetZapLog(zapcore.WarnLevel, 3)
	return &LogHelper{
		logger:  l,
		msgKey:  fmt.Sprintf("%s:nsq", defaultHelperLogger.msgKey),
		sprint:  defaultHelperLogger.sprint,
		sprintf: defaultHelperLogger.sprintf,
	}
}

// Log Print log by level and keyvals.
func Log(level log.Level, keyvals ...interface{}) error {
	_ = GetGlobalog().Log(level, keyvals...)
	return nil
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	defaultHelperLogger.Debug(a)
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelDebug, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	_ = GetGlobalog().Log(log.LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelInfo, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelInfo, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelWarn, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelWarn, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelError, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelError, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelFatal, defaultHelperLogger.msgKey, defaultHelperLogger.sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	_ = GetGlobalog().Log(log.LevelFatal, defaultHelperLogger.msgKey, defaultHelperLogger.sprintf(format, a...))
	os.Exit(1)
}
