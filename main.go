package main

import (
	"flag"
	"log"
	"runtime/debug"
)

const (
	nolimit = 0
)

var (
	Version   string
	BuildDate string
)

var (
	port   uint
	IPv6   bool
	verbos bool
	heap   int64
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.UintVar(&port, "port", 10086, "socks server listen port. ")
	flag.BoolVar(&IPv6, "ipv6", false, "enable ipv6. ")
	flag.BoolVar(&verbos, "verbos", false, "verbose log. ")
	flag.Int64Var(&heap, "memSize", nolimit, "memory size limit. ")
	flag.Parse()
}

func main() {
	if heap > 0 {
		debug.SetMemoryLimit(heap)
	}
	log.Printf("version: %s. build at: %s.", Version, BuildDate)
	exitChan := make(chan error)

	opts := make([]Opt, 0, 3)
	opts = append(opts, WithPort(uint16(port)))
	if verbos {
		opts = append(opts, EnableVerbos())
	}

	if IPv6 {
		opts = append(opts, EnableIpv6())
	}
	Start(exitChan, opts...)
	err := <-exitChan
	log.Fatalf("start with error %s", err.Error())
}
