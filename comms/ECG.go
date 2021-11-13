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
	ECGStreamPort        = "8884"
	MinDanceMovesToFatigue = 3
	ECGMarginToFatigue = 0.06
	ECGUserPort = 8882
)

var (
	ECGDataStream *ECGStream
	ECGOnce               sync.Once
)

type ECGStream struct {
	ecgStream chan session.ECG
	start          bool
	pastECGMax		uint32
	currentECGMax		uint32
}

func GetECGStream() *ECGStream {
	ECGOnce.Do(func() {
		ecgStream := make(chan session.ECG, bufferLength)
		ECGDataStream = &ECGStream{ecgStream: ecgStream,
			start: testPosition}
		go ECGDataStream.ClientListen()
	})

	return ECGDataStream
}

func (ecg *ECGStream) alertIfFatigued(ecgValue *session.ECG) {
	if ecg.start && ecgValue.Val3 > ecg.currentECGMax {
		ecg.currentECGMax = ecgValue.Val3
		fmt.Println("ECG current Max: ", ecg.currentECGMax, ", ECG past Max: ", ecg.pastECGMax)
	} else if !ecg.start && ecg.currentECGMax > 0 {
		ecg.pastECGMax = ecg.currentECGMax
		ecg.currentECGMax = 0
	}

	if float64(ecg.currentECGMax) > float64(ecg.pastECGMax) * (1 + ECGMarginToFatigue) &&
		len(GetStreamBuffer().PortMap[ECGUserPort]) > MinDanceMovesToFatigue { //trigger alert
		alert := &session.Alert{Message: "Muscles Fatigued! Do exercise caution!"}
		fmt.Println(">>>Muscles Fatigued!")

		ecg.pastECGMax = ecg.currentECGMax

		//might need mutex here to protect stream object todo
		GetStream(ECGUserPort).alert <- *alert
	}
}

func (ecg *ECGStream) ReadStream() session.ECG {
	ecgValue := <-ecg.ecgStream
	ecg.alertIfFatigued(&ecgValue)
	return ecgValue
}

func (ecg *ECGStream) Clear() {
	ecg.UpdateCommandStream(false)
	ecg.pastECGMax = 0
	ecg.currentECGMax = 0
}

func (ecg *ECGStream) UpdateCommandStream(state bool) {
	ecg.start = state
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

		//if ecg.start {
		//	GetStreamBuffer().Position = append(GetStreamBuffer().Position, *ecgData)
		//}
		ecg.ecgStream <- *ecgData
	}
}
