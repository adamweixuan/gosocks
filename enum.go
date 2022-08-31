package main

const (
	socksVersion = 0x05
)

const (
	defaultBufSize = 32 * 1024
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

type addrtype uint8

const (
	_      addrtype = iota
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

type statusCode uint8

const (
	_                             = iota
	success            statusCode = 0x00
	readConnErr        statusCode = 0x01
	connNotAllowed     statusCode = 0x02
	networkUnreachable statusCode = 0x03
	hostUnreachable    statusCode = 0x04
	connRefused        statusCode = 0x05
	ttlExpired         statusCode = 0x06
	cmdNotSupported    statusCode = 0x07
	addrNotSupported   statusCode = 0x08
)

type Message struct {
	code statusCode
	err  error
}

func NewMessage(code statusCode, err error) *Message {
	return &Message{code: code, err: err}
}

func (m *Message) Code() statusCode {
	return m.code
}

func (m *Message) SetCode(code statusCode) {
	m.code = code
}

func (m *Message) Err() error {
	return m.err
}

func (m *Message) SetErr(err error) {
	m.err = err
}

func NewReply(sc statusCode) []byte {
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
