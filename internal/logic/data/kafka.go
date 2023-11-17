package data

import (
	"github.com/2pgcn/gameim/conf"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func NewKafkaWriter(c *conf.Data_Kafka) (*kafka.Writer, error) {
	plainMechanism := plain.Mechanism{
		Username: c.GetSasl().GetUsername(),
		Password: c.GetSasl().GetPassword(),
	}
	return &kafka.Writer{
		Addr:         kafka.TCP(c.GetBootstrapServers()...),
		Topic:        c.GetTopic(),
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  10,
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		Transport: &kafka.Transport{
			SASL: plainMechanism,
		},
	}, nil
}
