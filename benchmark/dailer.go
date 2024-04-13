package main

import (
	"fmt"
	"net"
	"time"
)

// Dialer .
type Dialer struct {
	laddrIP string
	err     error
	dialer  *net.Dialer
}

// DialFromInterface .
func DialFromInterface(ifaceName string) *Dialer {
	d := &Dialer{}

	// Lookup rquested interface.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		d.err = err
		return d
	}

	// Pull the addresses.
	addres, err := iface.Addrs()
	if err != nil {
		d.err = err
		return d
	}

	// Look for the first usable address.
	var targetIP string
	for _, addr := range addres {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			d.err = err
			return d
		}
		if ip.IsUnspecified() {
			continue
		}
		if ip.To4().Equal(ip) {
			targetIP = ip.String()
		} else {
			targetIP = "[" + ip.String() + "]"
		}
	}
	if targetIP == "" {
		d.err = fmt.Errorf("no ipv4 found for interface")
		return d
	}
	d.laddrIP = targetIP
	return d
}

func (d *Dialer) lookupAddr(network, addr string) (net.Addr, error) {
	if d.err != nil {
		return nil, d.err
	}
	// If no custom dialer specified, use default one.
	if d.dialer == nil {
		d.dialer = &net.Dialer{}
	}

	// Resolve the address.
	switch network {
	case "tcp", "tcp4", "tcp6":
		addr, err := net.ResolveTCPAddr(network, d.laddrIP+":0")
		return addr, err
	case "udp", "udp4", "udp6":
		addr, err := net.ResolveUDPAddr(network, d.laddrIP+":0")
		return addr, err
	default:
		return nil, fmt.Errorf("unkown network")
	}
}

// Dial .
func (d *Dialer) Dial(network, addr string) (net.Conn, error) {
	laddr, err := d.lookupAddr(network, addr)
	if err != nil {
		return nil, err
	}
	d.dialer.LocalAddr = laddr
	return d.dialer.Dial(network, addr)
}

// DialTimeout .
func (d *Dialer) DialTimeout(network, addr string, timeout time.Duration) (net.Conn, error) {
	laddr, err := d.lookupAddr(network, addr)
	if err != nil {
		return nil, err
	}
	d.dialer.Timeout = timeout
	d.dialer.LocalAddr = laddr
	return d.dialer.Dial(network, addr)
}

// WithDialer .
func (d *Dialer) WithDialer(dialer net.Dialer) *Dialer {
	d.dialer = &dialer
	return d
}
