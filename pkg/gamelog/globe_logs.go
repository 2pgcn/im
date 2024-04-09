package gamelog

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
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
