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
	ECGStreamPort = "8884"
)

var (
	ECGDataStream *ECGStream
	ECGOnce               sync.Once
)

type ECGStream struct {
	ecgStream chan session.ECG
	commandStream  chan bool
	start          bool
}

func GetECGStream() *ECGStream {
	ECGOnce.Do(func() {
		ecgStream := make(chan session.ECG, bufferLength)
		commandStream := make(chan bool, bufferLength)
		ECGDataStream = &ECGStream{ecgStream: ecgStream,
			commandStream: commandStream, start: testPosition}
		go ECGDataStream.ClientListen()
	})

	return ECGDataStream
}

func (ecg *ECGStream) ReadStream() session.ECG {
	return <-ecg.ecgStream
}

func (ecg *ECGStream) UpdateCommandStream(state bool) {
	ecg.commandStream <- state
}

func (ecg *ECGStream) checkCommandStream() bool {
	if len(ecg.commandStream) == 0 {
		return false
	}

	return <-ecg.commandStream
}

func (ecg *ECGStream) ClientListen() {
	connPort := ECGStreamPort
	fmt.Println("Listening to Comms-Ultra96 via " + connType + " on " + connHost + ":" + connPort)
	c, err := net.Dial(connType, connHost+":"+connPort)
	for err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	defer c.Close()

	fmt.Println("Successfully Connected to Ultra96 on port:" + connPort)
	ecg.handleRequest(c)
}

func (ecg *ECGStream) handleRequest(conn net.Conn) {
	br := bufio.NewReader(conn)

	for {
		ecgRawData, err := br.ReadBytes('\x7F')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}

		ecgData := &session.ECG{}
		err = proto.Unmarshal(ecgRawData, ecgData)
		if err != nil {
			fmt.Println("Unmarshall Error", err.Error())
		}

		ecg.checkCommandStream()

		//if ecg.start {
		//	GetStreamBuffer().Position = append(GetStreamBuffer().Position, *ecgData)
		//}
		ecg.ecgStream <- *ecgData
	}
}
