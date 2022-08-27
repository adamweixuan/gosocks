package main

type Network uint8

const (
	_ Network = iota
	tcp
	udp
)

func (nt Network) String() string {
	switch nt {
	case tcp:
		return "tcp"
	case udp:
		return "udp"
	default:
		return "tcp"
	}
}

func NetworkFromStr(str string) Network {
	switch str {
	case "tcp", "tcp6":
		return tcp
	case "udp", "udp6":
		return udp
	default:
		return tcp
	}
}

type Addrtype uint8

const (
	_      Addrtype = iota
	ipv4            = 1
	domain          = 3
	ipv6            = 4
)

type IPVersion uint8

const (
	V4 IPVersion = iota
	V6
)
