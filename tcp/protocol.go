package tcp

import (
	"net"
)

type Packet interface {
	Format() string // data format
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
}
