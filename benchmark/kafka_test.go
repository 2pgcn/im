package main

import (
	"context"
	"github.com/2pgcn/gameim/internal/logic/data"
	"github.com/segmentio/kafka-go"
	"testing"
)

func TestKafka(t *testing.T) {
	conf := getLogicConfig()
	kafkaWr, err := data.NewKafkaWriter(conf.Data.Kafka)
	if err != nil {
		panic(err)
	}
	for i := 0; i <= 10; i++ {
		err = kafkaWr.WriteMessages(context.Background(), kafka.Message{
			Value: []byte("test 1111"),
		})
		if err != nil {
			t.Error(err)
		}
	}
}
