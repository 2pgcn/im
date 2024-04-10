package event

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"net"
	"time"
)

// 封包 简单的[len,data] len=4byte
// 为了本地测试
const defaultHeaderLen = 4

type SockRender struct {
	ctx    context.Context
	client net.Conn
	br     *bufio.Writer
}

func NewSockRender(addr string) (Sender, error) {
	var d net.Dialer
	raddr := net.UnixAddr{Name: addr, Net: "unix"}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := d.DialContext(ctx, "unix", raddr.String())
	if err != nil {
		return nil, err
	}
	return &SockRender{
		client: conn,
		br:     bufio.NewWriter(conn),
	}, nil
}

func (r *SockRender) Send(ctx context.Context, msg Event) error {
	//todo sendto 对应topic+appid
	m := msg.Value()
	data := make([]byte, defaultHeaderLen+len(m))
	binary.BigEndian.PutUint32(data, uint32(len(m)))
	copy(data[defaultHeaderLen:], m)
	err := binary.Write(r.client, binary.BigEndian, data)
	if err != nil {
		return err
	}
	return nil
}
func (r *SockRender) Close() error {
	return r.client.Close()
}

type SockReceiver struct {
	ctx       context.Context
	client    *net.UnixListener
	msgHandel *MessageHandler
	gopool    *safe.GoPool
}

func NewSockReceiver(con *conf.Sock) (*SockReceiver, error) {
	gopool := safe.NewGoPool(context.Background(), "QueueSockReceiver")
	raddr := net.UnixAddr{Name: con.GetAddress(), Net: "unix"}
	listenUnix, err := net.ListenUnix("unix", &raddr)
	if err != nil {
		return nil, err
	}
	msgHandel := &MessageHandler{
		queueNum:      defaultQueueNum,
		receiverQueue: []chan Event{},
	}
	for i := 0; i < msgHandel.queueNum; i++ {
		msgHandel.receiverQueue = append(msgHandel.receiverQueue, make(chan Event, defaultQueueLen))
	}
	r := &SockReceiver{
		client:    listenUnix,
		gopool:    gopool,
		msgHandel: msgHandel,
	}
	gopool.GoCtx(func(ctx context.Context) {
		//var conns []net.Conn
		index := 0
		for {
			conn, err := listenUnix.AcceptUnix()
			if err != nil {
				//todo,写入fin msg关闭conn,而不是关闭listenUnix退出
				gamelog.GetGlobalog().Errorf("queue sock accept error:%", err)
				return
			}
			gopool.GoCtx(func(ctx context.Context) {
				var head [defaultHeaderLen]byte
				redBuf := bufio.NewReader(conn)
				for {
					err := binary.Read(redBuf, binary.BigEndian, head[:])
					if err != nil {
						continue
					}
					packLen := binary.BigEndian.Uint32(head[:])
					if err != nil {
						gamelog.GetGlobalog().Errorf("queue sock read error:%err", err)
						continue
					}
					//读body
					body := make([]byte, packLen)
					err = binary.Read(redBuf, binary.BigEndian, &body)
					if err != nil {
						gamelog.GetGlobalog().Errorf("queue sock read error:%err", err)
						continue
					}
					msg := &QueueMsg{}
					err = json.Unmarshal(body, msg)
					if err != nil {
						gamelog.GetGlobalog().Errorf("queue sock read error:%err", err)
						continue
					}
					r.msgHandel.receiverQueue[index/defaultQueueNum] <- msg
				}
			})
		}
	})
	return r, nil
}

func (r *SockReceiver) Receive(ctx context.Context) (e []chan Event, err error) {
	return r.msgHandel.receiverQueue, nil
}

// 仅本地测试使用,未实现提交
func (r *SockReceiver) Commit(ctx context.Context, event Event) error {
	return nil
}

// 本地测试使用,未优雅退出,如处理掉SockReceiver里消息
func (r *SockReceiver) Close() error {
	return r.client.Close()
}
