package main

import (
	"bufio"
	"context"
	"github.com/2pgcn/gameim/api/protocol"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"
)

func main1() {
	cf, err := os.Create("/tmp/demo.cpu")
	if err != nil {
		panic(err)
	}
	defer func() {
		pprof.StopCPUProfile()
		cf.Close()
	}()
	pprof.StartCPUProfile(cf)
	ctx, cancel := context.WithCancel(context.Background())
	go func() { Server(ctx) }()
	time.Sleep(time.Second * 5)
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	br, bw := bufio.NewReader(conn), bufio.NewWriter(conn)
	msg := protocol.Proto{}
	msg.Op = protocol.OpSendMsg
	msg.Data = []byte("hello world")
	go func() {
		for {
			msg.WriteTcp(bw)
			addCountSend(1)
		}
	}()
	go func() {
		for {
			msg.DecodeFromBytes(br)
			addCountDown(1)
		}
	}()
	go func() {
		result(ctx)
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	select {
	case <-signals:
		cancel()
	}
}

func Server(ctx context.Context) {
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	listen, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.AcceptTCP()
		defer conn.Close()
		if err != nil {
			continue
		}
		go func() {
			for {
				br := bufio.NewReader(conn)
				msg := protocol.Proto{}
				msg.DecodeFromBytes(br)
			}
		}()
		go func() {
			for {
				bw := bufio.NewWriter(conn)
				msg := protocol.Proto{}
				msg.Op = protocol.OpSendMsg
				msg.Data = []byte("hello world")
				msg.WriteTcp(bw)
			}
		}()
	}
}

func startPprof() {
	cf, err := os.Create("/tmp/demo.cpu")
	if err != nil {
		panic(err)
	}
	defer func() {
		pprof.StopCPUProfile()
		cf.Close()
	}()
	pprof.StartCPUProfile(cf)
}
