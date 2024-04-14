package main

// Start Commond eg: ./client 1 1000 localhost:3101
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"time"

	"strconv"
	"sync/atomic"
)

var log *zap.SugaredLogger
var interfaceNames []string = []string{"eth0", "eth1", "eth2", "eth3", "eth4", "eth5"} // ens33、ens160

// only tcp and linux
func benchComet(ctx context.Context, addr string, num int) {
	flag.Parse()
	log = zap.NewExample().Sugar()
	begin := 1
	go result(ctx)
	for i := begin; i < begin+num; i++ {
		i = i
		go clients(ctx, addr, int64(i))
	}

}

func clients(ctx context.Context, addr string, mid int64) {
	startClient(ctx, addr, mid)
}

func startClient(ctx context.Context, addr string, key int64) {
	//time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
	atomic.AddInt64(&aliveCount, 1)
	defer atomic.AddInt64(&aliveCount, -1)
	// connnect to server
	index := rand.Intn(len(interfaceNames))
	ief, err := net.InterfaceByName(interfaceNames[index])
	if err != nil {
		log.Fatal(err)
	}
	addrs, err := ief.Addrs()
	if err != nil {
		log.Fatal(err)
	}
	//取最后一个,有ipv6地址
	laddr := &net.TCPAddr{
		IP: addrs[len(addrs)-1].(*net.IPNet).IP,
	}
	raddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", laddr, raddr)
	//defer conn.Close()
	//conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("net.Dial(%s) error(%v)", address, err)
		return
	}
	seq := uint16(0)
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)

	authToken := key
	//p := protocol.ProtoPool.Get()
	//defer protocol.ProtoPool.Put(p)
	p := &protocol.Proto{}
	p.Version = 1
	p.Op = protocol.OpAuth
	p.Seq = seq

	data, err := json.Marshal(&protocol.Auth{
		Appid: "app001",
		Token: strconv.Itoa(int(authToken)),
	})
	if err != nil {
		panic(err)
	}
	p.Data = data
	if err = p.WriteTcp(wr); err != nil {
		log.Errorf("tcpWriteProto() error(%v)", err)
		return
	}
	for {
		if err = p.DecodeFromBytes(rd); err == nil && p.Op == protocol.OpAuthReply {
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
	// heartbeat
	hbProto := &protocol.Proto{}
	hbProto.Op = protocol.OpSendMsg
	hbProto.Seq = seq
	hbProto.Data, _ = proto.Marshal(&protocol.Msg{
		Type: protocol.Type_PUSH,
		ToId: strconv.Itoa(int(authToken)),
		//SendId: userInfo.Uid,
		Msg: []byte("hello world gameim"),
	})

	go func() {
		for {
			time.Sleep(time.Second * 1)
			if err := hbProto.WriteTcp(wr); err != nil {
				log.Errorf("key:%d tcpWriteProto() error(%v)", key, err)
				return
			}
			addCountSend(1)
		}
	}()
	//log.Infof("key:%d send msg %+v", key, hbProto)
	seq++
	go func() {
		for {
			if err := p.DecodeFromBytes(rd); err != nil {
				gamelog.GetGlobalog().Error(err)
				return
			}
			addCountDown(1)
		}
	}()
}
