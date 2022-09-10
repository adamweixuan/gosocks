package main

import (
	"time"
)

type (
	Options struct {
		tcpNodelay bool
		reuseAddr  bool
		verbos     bool
		network    Network
		port       uint16
		timeout    time.Duration
		iface      string
	}
	Opt func(*Options)
)

func DefaultOpt() *Options {
	return &Options{
		port:       defaultPort,
		verbos:     false,
		iface:      "",
		network:    tcp,
		timeout:    time.Second,
		tcpNodelay: false,
	}
}

func WithPort(port uint16) Opt {
	return func(options *Options) {
		options.port = port
	}
}

func EnableVerbos() Opt {
	return func(options *Options) {
		options.verbos = true
	}
}

func EnableReuseAddr() Opt {
	return func(options *Options) {
		options.reuseAddr = true
	}
}

func WithTimeout(timeout time.Duration) Opt {
	return func(options *Options) {
		options.timeout = timeout
	}
}

func EnableNoDelay() Opt {
	return func(options *Options) {
		options.tcpNodelay = true
	}
}

func WithNetwork(nw string) Opt {
	return func(options *Options) {
		options.network = NetworkFromStr(nw)
	}
}

func WithIface(iface string) Opt {
	return func(options *Options) {
		options.iface = iface
	}
}
