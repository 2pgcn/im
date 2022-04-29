package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func main() {
	kconfig := sarama.NewConfig()
	kconfig.Version = sarama.V2_7_1_0 //
	kconfig.Consumer.Return.Errors = true
	kconfig.Consumer.Offsets.AutoCommit.Enable = false             // 自动提交
	kconfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
	kconfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kconfig.Consumer.Offsets.Retry.Max = 3
	kconfig.Consumer.MaxWaitTime = 1 * time.Second
	kconfig.Consumer.Fetch.Min = 10000
	kconfig.Consumer.Fetch.Default = 100000
	kconfig.Consumer.Fetch.Max = 1000000
	groupKafka, err := sarama.NewConsumerGroup([]string{"124.71.19.25:9094"}, "group", kconfig)
	if err != nil {
		panic(err)
	}
	defer func() { _ = groupKafka.Close() }()

	// Track errors
	go func() {
		for err := range groupKafka.Errors() {
			fmt.Println("ERROR", err)
		}
	}()
	consumerGroupTest := &consumerGroupTest{}
	for {
		err := groupKafka.Consume(context.Background(), []string{"comet0"}, consumerGroupTest)
		if err != nil {
			panic(err)
		}
	}

}

type consumerGroupTest struct {
}

func (consumer consumerGroupTest) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumer consumerGroupTest) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (consumer consumerGroupTest) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println(claim.Partition(), claim.Topic(), claim.InitialOffset(), claim.Messages())
	for msg := range claim.Messages() {
		fmt.Println(msg)
		sess.MarkMessage(msg, "")
		return nil
	}
	return nil
}
