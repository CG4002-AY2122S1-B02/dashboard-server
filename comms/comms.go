package comms

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

var packetStream chan Packet

func InitComm() {
	packetStream = make(chan Packet, 256)
}

func ReadStream() Packet {
	return <-packetStream
}

func ClientListen() {
	fmt.Println("Listening to Comms-Ultra96 via " + connType + " on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleRequest(c)
	}
}

func handleRequest(conn net.Conn) {
	br := bufio.NewReader(conn)

	for {
		packetData, err := br.ReadBytes('\x7F')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}

		packet := &Packet{}
		err = proto.Unmarshal(packetData, packet)
		if err != nil {
			fmt.Println("Unmarshall Error", err.Error())
		}

		fmt.Println(packet)
		packetStream <- *packet
	}
}