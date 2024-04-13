package main

import (
	"golang.org/x/sys/unix"
	"net"
	"syscall"
)

func setSocketOptions(network, address string, c syscall.RawConn, interfaceName string) (err error) {
	if interfaceName == "" || network != "tcp" {
		return
	}
	var innerErr error
	err = c.Control(func(fd uintptr) {
		host, _, _ := net.SplitHostPort(address)
		if ip := net.ParseIP(host); ip != nil && !ip.IsGlobalUnicast() {
			return
		}
		// 核心代理
		if interfaceName != "" {
			if innerErr = unix.BindToDevice(int(fd), interfaceName); innerErr != nil {
				return
			}
		}
	})
	if innerErr != nil {
		err = innerErr
	}
	return
}
