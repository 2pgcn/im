package event

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/trace_conf"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type kafkaSender struct {
	writer *kafka.Writer
	topic  string
}

func (s *kafkaSender) Send(ctx context.Context, message Event) error {
	ctx, span := trace_conf.SetTrace(ctx, trace_conf.LOGIC_PRODUCE_MSG, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()
	otel.GetTextMapPropagator().Inject(ctx, message.(*Msg))
	err := s.writer.WriteMessages(ctx, kafka.Message{
		Headers: message.Header().GetKafkaHead(),
		Key:     message.Key(),
		Value:   message.Value(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *kafkaSender) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaSender(c *conf.Data_Kafka) (Sender, error) {
	plainMechanism := plain.Mechanism{
		Username: c.GetSasl().GetUsername(),
		Password: c.GetSasl().GetPassword(),
	}
	w := &kafka.Writer{
		Addr:         kafka.TCP(c.GetBootstrapServers()...),
		Topic:        c.GetTopic(),
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  10,
		RequiredAcks: kafka.RequireAll,
		Async:        true,
		Transport: &kafka.Transport{

			SASL: plainMechanism,
		},
	}
	return &kafkaSender{writer: w, topic: c.GetTopic()}, nil
}

type kafkaReceiver struct {
	reader *kafka.Reader
	topic  string
}

func (k *kafkaReceiver) Receive(ctx context.Context) (e Event, err error) {
	//获取头数据
	m, err := k.reader.FetchMessage(ctx)
	if err != nil {
		return e, err
	}
	var msgData comet.MsgData
	err = proto.Unmarshal(m.Value, &msgData)
	if err != nil {
		return e, err
	}

	err = k.reader.CommitMessages(ctx, m)
	if err != nil {
		return e, err
	}
	header := NewHeader(8)
	for _, v := range m.Headers {
		header[v.Key] = string(v.Value)
	}
	header["partition"] = m.Partition
	//将topic,par,offset 保存到head处理完后提交
	e = &Msg{
		H:    header,
		K:    m.Key,
		Data: &msgData,
	}
	gamelog.Debug(m.Headers, e.Header())
	ctxTra := otel.GetTextMapPropagator().Extract(context.Background(), e.(*Msg))
	ctx, span := trace_conf.SetTrace(ctxTra, trace_conf.COMET_RECV_TO_QUEUE_MSG, trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()
	otel.GetTextMapPropagator().Inject(ctx, e.(*Msg))
	return e, err
}

func (k *kafkaReceiver) Close() error {
	err := k.reader.Close()
	if err != nil {
		return err
	}
	return nil
}

func (k *kafkaReceiver) Commit(ctx context.Context, event Event) error {
	msg, ok := event.(*Msg)
	if !ok {
		return fmt.Errorf("Commit msg error:%+v", event)
	}
	kMsg, err := msg.GetKafkaCommitMsg()
	if err != nil {
		return err
	}
	return k.reader.CommitMessages(ctx, kMsg)
}

func NewKafkaReceiver(ctx context.Context, c *conf.QueueMsg_Kafka) (Receiver, error) {
	kafkaReadConf := kafka.ReaderConfig{
		GroupID:  c.ConsumerGroupId,
		Brokers:  c.BootstrapServers,
		Topic:    c.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		//todo
		Logger:      gamelog.GetGlobeKafkaLog(),
		ErrorLogger: gamelog.GetGlobeKafkaErrorLog(),
	}
	if c.Sasl != nil {
		mechanism := plain.Mechanism{
			Username: c.Sasl.Username,
			Password: c.Sasl.Password,
		}
		dialer := &kafka.Dialer{
			Timeout:       3 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
			KeepAlive:     time.Second * 30,
		}
		kafkaReadConf.Dialer = dialer
	}
	kafkaRead := kafka.NewReader(kafkaReadConf)
	return &kafkaReceiver{reader: kafkaRead, topic: c.Topic}, nil
}
