package main

import (
	"flag"
	"math"
	"runtime/debug"
	"time"
)

var (
	Version   string
	BuildDate string
)

var (
	port      uint
	verbos    bool
	heap      int64
	iface     string
	reuseAddr bool
	timeout   time.Duration
	noDelay   bool
	network   string
)

func init() {
	flag.UintVar(&port, "port", 10086, "set socks server listen port. ")
	flag.BoolVar(&verbos, "verbos", false, "enable verbose log. ")
	flag.Int64Var(&heap, "heap", math.MaxInt64, "set memory size limit. ")
	flag.StringVar(&iface, "iface", "", "set specified interface to use. ")
	flag.BoolVar(&reuseAddr, "reuseaddr", true, "enable reuseaddr. ")
	flag.DurationVar(&timeout, "timeout", time.Second, "set session timeout. ")
	flag.BoolVar(&noDelay, "nodelay", true, "enable tcpnodelay. ")
	flag.StringVar(&network, "network", "tcp", "set network tcp or udp.")
	flag.Parse()
}

func parseOpts() []Opt {
	opts := make([]Opt, 0, 7)
	opts = append(opts, WithPort(uint16(port)))
	opts = append(opts, WithNetwork(network))
	opts = append(opts, WithTimeout(timeout))

	if verbos {
		opts = append(opts, EnableVerbos())
	}

	if reuseAddr {
		opts = append(opts, EnableReuseAddr())
	}

	if noDelay {
		opts = append(opts, EnableNoDelay())
	}

	if len(iface) > 0 {
		opts = append(opts, WithIface(iface))
	}
	return opts
}

func main() {
	debug.SetMemoryLimit(heap)
	log.Info("version: %s. build at: %s.", Version, BuildDate)
	exitChan := make(chan error)

	Start(exitChan, parseOpts()...)

	err := <-exitChan
	log.Error("start with error %s", err.Error())
}
