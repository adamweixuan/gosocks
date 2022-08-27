package main

import (
	"fmt"
	"net"
)

func Start(exitChan chan error, opts ...Opt) {
	opt := DefaultOpt()
	for _, op := range opts {
		op(opt)
	}
	go func() {
		if err := run(opt); err != nil {
			exitChan <- err
		}
	}()
}

func run(opt *Options) error {
	InitLocalIPByInterName(opt.iface)
	localIP, netIface := GetNetInterfaceByName(opt.iface)

	addr := fmt.Sprintf(":%d", opt.port)
	if localIP != nil {
		addr = fmt.Sprintf("%s:%d", localIP.String(), opt.port)
	}
	ln, err := net.Listen(opt.network.String(), addr)
	if err != nil {
		return err
	}
	if ln == nil {
		return ErrNilListener
	}

	log.Info("iface name: %v, out ip: %v. server start at: %s ", netIface, localIP.String(), ln.Addr().String())
	for {
		stream, err := ln.Accept()
		if err != nil {
			log.Error("Accept with error:%s", err.Error())
			continue
		}

		go func() {
			ctx := NewCtxWithTraceID()
			session := &Session{
				verbos:    opt.verbos,
				local:     stream,
				remote:    nil,
				ip:        localIP,
				iface:     netIface,
				ntype:     opt.network,
				timeout:   opt.timeout,
				noDelay:   opt.tcpNodelay,
				reuseAddr: opt.reuseAddr,
			}
			session.Start(ctx)
		}()
	}
}
