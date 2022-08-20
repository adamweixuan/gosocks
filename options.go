package main

type IPVersion uint8

const (
	V4 IPVersion = iota
	V6
)

type (
	Options struct {
		port   uint16
		ipVer  IPVersion
		verbos bool
	}
	Opt func(*Options)
)

func DefaultOpt() *Options {
	return &Options{
		port:   10086,
		ipVer:  V4,
		verbos: false,
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

func EnableIpv6() Opt {
	return func(options *Options) {
		options.ipVer = V6
	}
}
