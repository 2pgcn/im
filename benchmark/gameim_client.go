package main

// Start Commond eg: ./client 1 1000 localhost:3101
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/php403/gameim/api/comet"
	"github.com/php403/gameim/api/logic"
	"github.com/php403/gameim/api/protocol"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

var log *zap.SugaredLogger

const (
	OpAuth      = uint16(2)
	OpAuthReply = uint16(3)

	OpSendMsg      = uint16(4)
	OpSendMsgReply = uint16(5)

	OpSendRoomMsg      = uint16(6)
	OpSendMsgRoomReply = uint16(7)

	OpSendAreaMsg      = uint16(8)
	OpSendAreaMsgReply = uint16(9)

	OpDisconnect      = uint16(10)
	OpDisconnectReply = uint16(11)
)

const (
	rawHeaderLen = uint16(16)
)

// AuthToken auth token.
type AuthToken struct {
	Mid      int64   `json:"mid"`
	Key      string  `json:"key"`
	RoomID   string  `json:"room_id"`
	Platform string  `json:"platform"`
	Accepts  []int32 `json:"accepts"`
}

var (
	countDown  int64
	aliveCount int64
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	log = zap.NewExample().Sugar()
	begin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	go result()
	for i := begin; i < begin+num; i++ {
		go client(int64(i))
	}
	// signal
	var exit chan bool
	<-exit
}

func result() {
	var (
		lastTimes int64
		interval  = int64(5)
	)
	for {
		nowCount := atomic.LoadInt64(&countDown)
		nowAlive := atomic.LoadInt64(&aliveCount)
		diff := nowCount - lastTimes
		lastTimes = nowCount
		fmt.Println(fmt.Sprintf("%s alive:%d down:%d down/s:%d", time.Now().Format("2006-01-02 15:04:05"), nowAlive, nowCount, diff/interval))
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func client(mid int64) {
	for {
		startClient(mid)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func startClient(key int64) {
	//time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
	atomic.AddInt64(&aliveCount, 1)
	quit := make(chan bool, 1)
	defer func() {
		close(quit)
		atomic.AddInt64(&aliveCount, -1)
	}()
	// connnect to server
	fmt.Println(os.Args[3])
	conn, err := net.Dial("tcp", os.Args[3])
	if err != nil {
		log.Errorf("net.Dial(%s) error(%v)", os.Args[3], err)
		return
	}
	seq := uint16(0)
	wr := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)

	p := new(protocol.Proto)
	p.Version = 1
	p.Op = OpAuth
	p.Seq = seq
	p.Data = []byte("testToken")
	if err = p.WriteTcp(wr); err != nil {
		log.Errorf("tcpWriteProto() error(%v)", err)
		return
	}
	for {
		if err = p.DecodeFromBytes(rd); err == nil && p.Op == OpAuthReply {
			log.Infof("key:%d auth ok, p: %v", key, p)
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
	//测试 1000000用户在一个区
	go func() {
		hbProto := new(protocol.Proto)
		for {
			// heartbeat
			hbProto.Op = OpSendAreaMsg
			hbProto.Seq = seq
			hbProto.Data, _ = proto.Marshal(&comet.Msg{
				Type:   comet.Type_AREA,
				ToId:   userInfo.AreaId,
				SendId: userInfo.Uid,
				Msg:    []byte("hello world"),
			})
			if err = hbProto.WriteTcp(wr); err != nil {
				log.Errorf("key:%d tcpWriteProto() error(%v)", key, err)
				return
			}
			//log.Infof("key:%d send msg %+v", key, hbProto)
			seq++
			select {
			case <-quit:
				return
			default:
			}
			time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))
		}
	}()
	// reader
	for {
		if err = p.DecodeFromBytes(rd); err == nil {
			atomic.AddInt64(&countDown, 1)
			//log.Infof("key:%d auth ok, p: %v", key, p)
		}
	}
}
