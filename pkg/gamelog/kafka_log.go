package gamelog

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"math/rand"
	"sync"
)

var _ kafka.Logger = (*KafkaLog)(nil)
var one = sync.Once{}

var globeKafkaLog *KafkaLog
var globeKafkaErrorLog *KafkaLog

func GetGlobeKafkaLog() *KafkaLog {
	return globeKafkaLog
}

func GetGlobeKafkaErrorLog() *KafkaLog {
	return globeKafkaErrorLog
}

type KafkaLog struct {
	level log.Level
	l     GameLog
}

func (l *KafkaLog) Printf(msg string, args ...interface{}) {
	if l.level > log.LevelWarn {
		l.l.Errorf(msg, args)
		return
	}
	//debug日志太多,过滤10%
	if rand.Intn(10)%10 == 0 {
		l.l.Debug(msg, args)
	}
}

type KafkaErrorLog struct {
	l GameLog
}

func InitKafkaLog(l *LogHelper) {
	one.Do(func() {
		l = l.ReplacePrefix(DefaultMessageKey + "kafka").(*LogHelper)
		globeKafkaLog = &KafkaLog{l: l, level: log.LevelDebug}
		globeKafkaErrorLog = &KafkaLog{l: l, level: log.LevelError}
	})
}
