package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bmbstack/gosocket/example"
	"github.com/bmbstack/gosocket/tcp"
)

type Callback struct{}

func (this *Callback) OnConnect(c *tcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (this *Callback) OnMessage(c *tcp.Conn, p tcp.Packet) bool {
	packet := p.(*example.OnePacket)
	log.Println(fmt.Sprintf("FromClient: body.length=%d, format=%s, body=%v", packet.GetBodyLength(), packet.Format(), string(packet.GetBody())))
	c.AsyncWritePacket(example.NewOnePacket("text", []byte("tom")), time.Second)
	return true
}

func (this *Callback) OnClose(c *tcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	config := &tcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := tcp.NewServer(config, &Callback{}, &example.OneProtocol{})

	// starts service
	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
