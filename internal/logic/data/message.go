package data

import (
	"context"
	"github.com/2pgcn/gameim/pkg/event"
)

func (d *Data) WriteKafkaMessage(ctx context.Context, e event.Event) error {
	return d.kafkaClient.Send(ctx, e)
}
