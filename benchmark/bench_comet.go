package main

// Start Commond eg: ./client 1 1000 localhost:3101
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"bufio"
	"context"
	"flag"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var log *zap.SugaredLogger

func benchComet(ctx context.Context, addr string, num int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	log = zap.NewExample().Sugar()
	begin := 0
	log.Debug(num)
	go result(ctx)
	for i := begin; i < begin+num; i++ {
		go clients(ctx, addr, int64(i))
	}

}

func clients(ctx context.Context, addr string, mid int64) {
	for {
		startClient(ctx, addr, mid)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func startClient(ctx context.Context, addr string, key int64) {
	//time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
	atomic.AddInt64(&aliveCount, 1)
	defer atomic.AddInt64(&aliveCount, -1)
	// connnect to server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("net.Dial(%s) error(%v)", address, err)
		return
	}
	seq := uint16(0)
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)

	authToken := atomic.LoadInt64(&aliveCount)
	p := protocol.ProtoPool.Get()
	defer protocol.ProtoPool.Put(p)
	p.Version = 1
	p.Op = protocol.OpAuth
	p.Seq = seq
	p.Data = []byte(strconv.Itoa(int(authToken)))
	log.Infof("auth start")
	if err = p.WriteTcp(wr); err != nil {
		log.Errorf("tcpWriteProto() error(%v)", err)
		return
	}
	log.Infof("auth success")
	for {
		if err = p.DecodeFromBytes(rd); err == nil && p.Op == protocol.OpAuthReply {
			log.Infof("key:%d auth ok, p: %v", strconv.FormatInt(key, 10), p)
			break
		}
	}
	seq++
	// writer
	userInfo := &logic.AuthReply{}
	err = proto.Unmarshal(p.Data, userInfo)
	if err != nil {
		log.Errorf("auth proto.Unmarshal() error(%v)", err)
	}
	addAliveCount(1)
	//测试 1000000用户在一个区
	//go func() {
	//hbProto := protocol.ProtoPool.Get()
	//for {
	// heartbeat
	hbProto := &protocol.Proto{}
	hbProto.Op = protocol.OpSendAreaMsg
	hbProto.Seq = seq
	hbProto.Data, _ = proto.Marshal(&protocol.Msg{
		Type:   protocol.Type_ROOM,
		ToId:   userInfo.RoomId,
		SendId: userInfo.Uid,
		Msg:    []byte("hello world gameim"),
	})
	go func() {
		for {
			if err = hbProto.WriteTcp(wr); err != nil {
				log.Errorf("key:%d tcpWriteProto() error(%v)", key, err)
				return
			}
			addCountSend(1)
			time.Sleep(time.Second * 1)
		}
	}()
	//log.Infof("key:%d send msg %+v", key, hbProto)
	seq++
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err = p.DecodeFromBytes(rd); err == nil {
				atomic.AddInt64(&countDown, 1)
				addCountDown(1)
				log.Infof("key:%v. p: %s", key, string(p.Data))
			}
		}
	}
}
