package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/php403/gameim/api/comet"
	pb "github.com/php403/gameim/api/logic"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync/atomic"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedLogicServer
	kafka sarama.AsyncProducer
}

var (
	port   = flag.Int("port", 9000, "The server port")
	online = uint64(1)
)

// OnAuth 100万连接 1个区分10个工会
func (s *server) OnAuth(ctx context.Context, in *pb.AuthReq) (*pb.AuthReply, error) {
	//token
	uid := atomic.AddUint64(&online, 1)
	return &pb.AuthReply{
		Uid:    uid,
		AreaId: uint64(1),
		RoomId: uid % 10000,
	}, nil
}

func (s *server) OnMessage(ctx context.Context, req *pb.MessageReq) (*empty.Empty, error) {
	//落盘或者鉴黄,或者校验是否数据不一致
	kafkaMeta := &comet.Msg{
		Type: req.Type,
		ToId: req.ToId,
		Msg:  req.Msg,
	}
	kafkaMetaByte, err := proto.Marshal(kafkaMeta)
	if err != nil {
		return &empty.Empty{}, err
	}
	s.kafka.Input() <- &sarama.ProducerMessage{
		Topic: req.GetCometKey(),
		Value: sarama.ByteEncoder(kafkaMetaByte),
	}
	return &empty.Empty{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	kafka := NewKafkaProducer()
	pb.RegisterLogicServer(s, &server{kafka: kafka})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
