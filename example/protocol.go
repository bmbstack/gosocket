package example

import (
	"encoding/binary"
	"errors"
	"github.com/bmbstack/gosocket/tcp"
	"io"
	"net"
)

const LengthSize = 4  // length size, 4 bytes, uint32
const FormatSize = 16 // format size, 16 bytes, string

type OnePacket struct {
	format string // text, image
	buff   []byte // buff = lengthBytes + formatBytes + body, header = lengthBytes + formatBytes
}

func NewOnePacket(format string, body []byte) *OnePacket {
	p := &OnePacket{}
	p.format = format
	p.buff = make([]byte, LengthSize+FormatSize+len(body))
	binary.BigEndian.PutUint32(p.buff[0:LengthSize], uint32(len(body)))
	copy(p.buff[LengthSize:LengthSize+FormatSize], format)
	copy(p.buff[LengthSize+FormatSize:], body)
	return p
}

func (this *OnePacket) Format() string {
	return this.format
}

func (this *OnePacket) Serialize() []byte {
	return this.buff
}

func (this *OnePacket) GetBodyLength() uint32 {
	return binary.BigEndian.Uint32(this.buff[0:LengthSize])
}

func (this *OnePacket) GetBody() []byte {
	return this.buff[LengthSize+FormatSize:]
}

type OneProtocol struct {
}

func (this *OneProtocol) ReadPacket(conn *net.TCPConn) (tcp.Packet, error) {
	var headerBytes = make([]byte, LengthSize+FormatSize)
	var length uint32

	// read header bytes (lengthBytes + formatBytes)
	if _, err := io.ReadFull(conn, headerBytes); err != nil {
		return nil, err
	}
	lengthBytes := headerBytes[0:LengthSize]
	formatBytes := headerBytes[LengthSize : LengthSize+FormatSize]

	length = binary.BigEndian.Uint32(lengthBytes)
	format := string(formatBytes)

	if format == "text" {
		if length > 4096 {
			return nil, errors.New("the size of packet is larger than the limit 4k")
		}
	} else if format == "image" {
		if length > 10240 {
			return nil, errors.New("the size of packet is larger than the limit 10k")
		}
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, err
	}

	return NewOnePacket(format, body), nil
}
