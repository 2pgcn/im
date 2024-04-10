package data

import (
	"context"
	"github.com/2pgcn/gameim/pkg/event"
)

func (d *Data) WriteKafkaMessage(ctx context.Context, e event.Event) error {
	//return d.kafkaClient.Send(ctx, e)
	return nil
}

func (d *Data) WriteMessage(ctx context.Context, e event.Event) error {
	//todo 根据appid send发送到不同topic
	return d.producer.Send(ctx, e)
}
