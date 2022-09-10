package main

const (
	socksVersion   = 0x05
	defaultBufSize = 32 * 1024
	defaultPort    = 10086
)

const (
	tcpStr  = "tcp"
	tcp6Str = "tcp6"
	udpStr  = "udp"
	udp6Str = "udp6"
)

type Network uint8

const (
	_ Network = iota
	tcp
	udp
)

func (nt Network) String() string {
	switch nt {
	case tcp:
		return tcpStr
	case udp:
		return udpStr
	default:
		return tcpStr
	}
}

func NetworkFromStr(str string) Network {
	switch str {
	case tcpStr, tcp6Str:
		return tcp
	case udpStr, udp6Str:
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

type cmdType byte

const (
	_                       = iota
	connectCmd      cmdType = 1
	bindCmd         cmdType = 2
	udpAssociateCmd cmdType = 3
)

type StatusCode uint8

const (
	success StatusCode = iota
	readConnErr
	connNotAllowed
	networkUnreachable
	hostUnreachable
	connRefused
	ttlExpired
	cmdNotSupported
	addrNotSupported
)

type Message struct {
	code StatusCode
	err  error
}

func NewMessage(code StatusCode, err error) *Message {
	return &Message{code: code, err: err}
}

func (m *Message) Code() StatusCode {
	return m.code
}

func (m *Message) SetCode(code StatusCode) {
	m.code = code
}

func (m *Message) Err() error {
	return m.err
}

func (m *Message) SetErr(err error) {
	m.err = err
}

func NewReply(sc StatusCode) []byte {
	return []byte{
		socksVersion,
		byte(sc),
		0x00,
		0x01,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}
}

const (
	noAuthRequired   byte = iota
	noAcceptableAuth byte = 255
)
