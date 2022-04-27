package main

import (
	"github.com/Shopify/sarama"
)

func NewKafkaProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.NoResponse
	producer, err := sarama.NewAsyncProducer([]string{"124.71.19.25:9094"}, config)
	if err != nil {
		panic(err)
	}
	return producer
}
