package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/bmbstack/gosocket/example"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	protocol := &example.OneProtocol{}

	// ping <--> pong
	for i := 1; i <= 100; i++ {
		// write
		conn.Write(example.NewOnePacket("text", []byte(fmt.Sprintf("hello%d", i))).Serialize())

		// read
		p, err := protocol.ReadPacket(conn)
		if err == nil {
			packet := p.(*example.OnePacket)
			log.Println(fmt.Sprintf("FromServer: body.length=%d, format=%s, body=%v", packet.GetBodyLength(), packet.Format(), string(packet.GetBody())))
		}

		time.Sleep(2 * time.Second)
	}

	_ = conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
