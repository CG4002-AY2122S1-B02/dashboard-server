package comms

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func ClientListen() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
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

		handleRequest(c)
	}
}

func handleRequest(conn net.Conn) {
	// Buffer that holds incoming information
	buf := make([]byte, 8)

	for {
		len, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		s := string(buf[:len])

		fmt.Println("Stuff", s)
		fmt.Println("len", binary.Size(buf))
	}
}