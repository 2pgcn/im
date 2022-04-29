package comet

import (
	"github.com/Shopify/sarama"
	"github.com/php403/gameim/config"
	"time"
)

// NewConsumerGroupHandler NewKafkaConsumer todo  NewKafkaConsumer改成接口 后续扩展
func NewConsumerGroupHandler(c *config.CometConfigQueueMsg) (consumer sarama.ConsumerGroup, err error) {
	kconfig := sarama.NewConfig()
	kconfig.Version = sarama.V2_7_1_0 //
	//kconfig.Consumer.Return.Errors = true
	kconfig.Consumer.Offsets.AutoCommit.Enable = true              // 自动提交
	kconfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
	kconfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kconfig.Consumer.Offsets.Retry.Max = 3
	kconfig.Consumer.MaxWaitTime = 1 * time.Second
	return sarama.NewConsumerGroup(c.Brokers, c.Group, kconfig)
}
