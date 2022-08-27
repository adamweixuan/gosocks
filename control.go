package main

import (
	"net"
	"syscall"
)

type CtrlOpt struct {
	iface      *net.Interface
	reuseAddr  bool
	tcpNodelay bool
}

type CtrlOp func(args *CtrlOpt)

func Bind(intf *net.Interface) CtrlOp {
	return func(opts *CtrlOpt) {
		opts.iface = intf
	}
}

func ReuseAddr() CtrlOp {
	return func(opts *CtrlOpt) {
		opts.reuseAddr = true
	}
}

func NoDelay() CtrlOp {
	return func(opts *CtrlOpt) {
		opts.tcpNodelay = true
	}
}

func Control(opts ...CtrlOp) func(network, address string, c syscall.RawConn) error {
	option := &CtrlOpt{}
	for _, opt := range opts {
		opt(option)
	}
	return control(option)
}
