package comms

import (
	"bufio"
	"dashboard-server/internal/session"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"sync"
)

const (
	positionStreamPort = "8880"
)

var (
	PositionDataStream *PositionStream
	once               sync.Once
)

type PositionStream struct {
	positionStream chan session.Position
	commandStream  chan bool
	start          bool
}

func GetPositionStream() *PositionStream {
	once.Do(func() {
		positionStream := make(chan session.Position, bufferLength)
		commandStream := make(chan bool, bufferLength)
		PositionDataStream = &PositionStream{positionStream: positionStream,
			commandStream: commandStream, start: testPosition}
		go PositionDataStream.ClientListen()
	})

	return PositionDataStream
}

func (p *PositionStream) ReadStream() session.Position {
	return <-p.positionStream
}

func (p *PositionStream) UpdateCommandStream(state bool) {
	p.commandStream <- state
}

func (p *PositionStream) checkCommandStream() bool {
	if len(p.commandStream) == 0 {
		return false
	}

	return <-p.commandStream
}

func (p *PositionStream) ClientListen() {
	connPort := positionStreamPort
	fmt.Println("Listening to Comms-Ultra96 via " + connType + " on " + connHost + ":" + connPort + " (position)")
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		log.Fatal("Error listening:", err.Error())
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

		go p.handleRequest(c)
	}
}

func (p *PositionStream) handleRequest(conn net.Conn) {
	br := bufio.NewReader(conn)

	for {
		positionData, err := br.ReadBytes('\x7F')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}

		position := &session.Position{}
		err = proto.Unmarshal(positionData, position)
		if err != nil {
			fmt.Println("Unmarshall Error", err.Error())
		}

		p.checkCommandStream()

		if p.start {
			GetStreamBuffer().Position = append(GetStreamBuffer().Position, *position)
		}
		p.positionStream <- *position
	}
}
