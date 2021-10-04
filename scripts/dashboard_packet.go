package main

import (
	"bufio"
	"dashboard-server/internal/session"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"os"
)

const (
	connHost     = "localhost"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	PORT := arguments[1]
	c, err := net.Dial("tcp", connHost+":"+ PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		//message, _ := bufio.NewReader(c).ReadString('\n')
		//fmt.Print("->: " + message)
		br := bufio.NewReader(c)
		packetData, err := br.ReadBytes('\x7F')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}
		if PORT == "8880" {
			packet := &session.Position{}
			err = proto.Unmarshal(packetData, packet)
			if err != nil {
				fmt.Println("Unmarshall Error", err.Error())
			}

			fmt.Println(packet)
		} else {
			packet := &session.Packet{}
			err = proto.Unmarshal(packetData, packet)
			if err != nil {
				fmt.Println("Unmarshall Error", err.Error())
			}

			fmt.Println(packet)
		}

	}
}
