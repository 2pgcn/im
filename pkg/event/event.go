package event

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type EventHeader map[string]any

type Event interface {
	Header() EventHeader
	Key() []byte
	Value() []byte
	RawValue() any
	StartTrace(traceName string)
	String() string
}

func NewHeader(size int) EventHeader {
	return make(EventHeader, size)
}

// todo 改成proto.Marshal
func (eh EventHeader) GetKafkaHead() (res []kafka.Header) {
	for k, v := range eh {
		res = append(res, kafka.Header{
			Key:   k,
			Value: []byte(v.(string)),
		})
	}
	return res
}

type Handler func(context.Context, Event) error

type Sender interface {
	Send(ctx context.Context, msg Event) error
	Close() error
}

type Receiver interface {
	Receive(ctx context.Context) (e Event, err error)
	Commit(ctx context.Context, event Event) error
	Close() error
}
