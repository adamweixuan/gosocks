package main

import (
	"errors"
	"fmt"
	"log"
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
	network := "tcp"
	if opt.ipVer == V6 {
		network = "tcp6"
	}
	addr := fmt.Sprintf(":%d", opt.port)
	ln, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	if ln == nil {
		return errors.New("listen fail")
	}

	sessionID := uint64(0)

	log.Printf("server start at: %s", ln.Addr().String())
	for {
		stream, err := ln.Accept()
		if err != nil {
			log.Printf("Accept with error:%s", err.Error())
			continue
		}
		sessionID++
		go func() {
			session := &Session{
				verbos: opt.verbos,
				local:  stream,
				id:     sessionID,
			}
			session.Start()
		}()
	}
}
