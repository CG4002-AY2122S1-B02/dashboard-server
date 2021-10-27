package comms

import (
	"bufio"
	"dashboard-server/internal/session"
	"dashboard-server/internal/stream/po"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"math"
	"net"
	"sync"
)

const (
	connHost     = "localhost"
	connType     = "tcp"
	bufferLength = 256
	testComms    = false
	testPosition = true

	/**
	accuracy
	~0.58
	~0.71
	~(8s) 0.50 dab
	~(9s) 0.675 jamesbond
	~0.533 mermaid
	~0.583 Dab
	~(9s) 0.756 jamesbond
	~(9s) 0.49 Mermaid
	 */

	lowerBound3star = 0.85
	lowerBound2star = 0.75
)

var (
	streamMap        map[int]*Stream
	streamBuffer     *StreamBuffer
	streamBufferOnce sync.Once
	portMapMutex = &sync.Mutex{}
	portStartMutex = &sync.Mutex{}
)

type StreamCommand struct {
	start            bool
	username         string
	accountName      string
	sessionTimestamp uint64
}

type StreamBuffer struct {
	PortMap  map[int][]session.Packet
	PortStatus  map[int]uint64
	Position       []session.Position
	groupSyncDelay chan uint64
	pointer        int
	TotalSyncDelay uint64
}

func (sb *StreamBuffer) Clear() {
	sb.PortMap = make(map[int][]session.Packet)
	sb.PortStatus = make(map[int]uint64)
	sb.Position = make([]session.Position, 0)
	sb.groupSyncDelay = make(chan uint64, bufferLength)
	sb.pointer = 0
	sb.TotalSyncDelay = 0
}

func (sb *StreamBuffer) ReadGroupSyncDelay() uint64 {
	return <- sb.groupSyncDelay
}

func (sb *StreamBuffer) UpdateStartGroupSyncDelay(packet session.Packet, port int) {
	//input all kinds of inputs, predicted, start, end. publish everytime start is received. Reset when 3 ends or predictions is received
	//if session packet buffer is full for all users at that index, can compute group sync delay

	//status > 0: start epoch, status = 0 or !ok: cleared awaiting new start, port*2 ==1: ended/predicted

	portStartMutex.Lock()
	if packet.DanceMove == "START" {
		sb.PortStatus[port] = packet.EpochMs
		//sb.PortStatus[port*2] = 0
	} else {
		sb.PortStatus[port*2] = 1
	}
	
	//if all packets have been predicted/ended, save to TotalSyncDelay and clear
	if sb.PortStatus[8881*2] == 1 && sb.PortStatus[8882*2] == 1 && sb.PortStatus[8883*2] == 1 {
		syncDelay := po.ComputeSyncDelay(
			[]uint64{sb.PortStatus[8881],sb.PortStatus[8882],sb.PortStatus[8883]},
		)

		if packet.Accuracy >= 0 || packet.DanceMove == "END" { //only add to totalSyncDelay if a prediction/end is passed. If reconnecting, result is discarded
			sb.TotalSyncDelay += syncDelay
			if syncDelay > 0 {
				sb.pointer += 1
			}
		}

		sb.PortStatus[8881*2] = 0
		sb.PortStatus[8882*2] = 0
		sb.PortStatus[8883*2] = 0
		sb.PortStatus[8881] = 0
		sb.PortStatus[8882] = 0
		sb.PortStatus[8883] = 0
	}
	portStartMutex.Unlock()

	status1, ok1 := sb.PortStatus[8881]
	status2, ok2 := sb.PortStatus[8882]
	status3, ok3 := sb.PortStatus[8883]
	
	//Publish if all ports have a start epoch
	if ok1 && status1 > 0 && ok2 && status2 > 0 && ok3 && status3 > 0 && packet.DanceMove == "START"{
		syncDelay := po.ComputeSyncDelay(
			[]uint64{status1,status2,status3},
			)
		sb.groupSyncDelay <- syncDelay
	}
}

func (sb *StreamBuffer) UpdateGroupSyncDelay() {
		//if session packet buffer is full for all users at that index, can compute group sync delay
	if len(sb.PortMap[8881]) < sb.pointer + 1 {
		//if buffer v has not reached next dance move, we cannot compute
		return
	}
	if len(sb.PortMap[8882]) < sb.pointer + 1 {
		//if buffer v has not reached next dance move, we cannot compute
		return
	}
	if len(sb.PortMap[8883]) < sb.pointer + 1 {
		//if buffer v has not reached next dance move, we cannot compute
		return
	}


	//all buffers would be at least len(sb.PortMap) long
	//now pointer points to the buffer slot to compute group sync delay

	syncDelay := po.ComputeSyncDelay(
		[]uint64{sb.PortMap[8881][sb.pointer].EpochMs,
			sb.PortMap[8882][sb.pointer].EpochMs,
			sb.PortMap[8883][sb.pointer].EpochMs})

	sb.groupSyncDelay <- syncDelay
	sb.TotalSyncDelay += syncDelay

	sb.pointer += 1
}

func (sb *StreamBuffer) GetAvgSyncDelay() uint64 {
	if sb.pointer == 0 {
		return 0
	}

	return uint64(math.Round(float64(sb.TotalSyncDelay) / float64(sb.pointer)))
}

func GetStreamBuffer() *StreamBuffer {
	streamBufferOnce.Do(func() {
		streamBuffer = &StreamBuffer{
			make(map[int][]session.Packet),
			make(map[int]uint64),
			make([]session.Position, 0),
			make(chan uint64, bufferLength),
			0,
			0,
		}

		streamBuffer.PortMap[8881] = make([]session.Packet, 0)
		streamBuffer.PortMap[8882] = make([]session.Packet, 0)
		streamBuffer.PortMap[8883] = make([]session.Packet, 0)
	})

	return streamBuffer
}

type Stream struct {
	port             int
	packetStream     chan session.Packet
	commandStream    chan StreamCommand
	sessionTimestamp uint64
	username         string
	accountName      string
	start            bool
	lastMove         chan string
	lastAccuracy     chan float32
	alert			chan session.Alert
}

func GetStream(port int) *Stream {
	return streamMap[port]
}

func NewStream(port int) *Stream {
	return InitialiseStream(port, testComms)
}

func InitialiseStream(port int,
	start bool) *Stream {
	if streamMap == nil {
		streamMap = make(map[int]*Stream)
	}

	streamMap[port] = &Stream{port: port,
		packetStream:  make(chan session.Packet, bufferLength),
		commandStream: make(chan StreamCommand, bufferLength),
		start:         start,
		lastMove:      make(chan string, bufferLength),
		lastAccuracy:  make(chan float32, bufferLength),
		alert: make(chan session.Alert, bufferLength),
	}

	go streamMap[port].ClientListen()
	return streamMap[port]
}

func (s *Stream) ReadStream() session.Packet {
	packet := <-s.packetStream
	if packet.Accuracy >= 0 {
		s.lastMove <- packet.DanceMove
		s.lastAccuracy <- packet.Accuracy
	}

	return packet
}

func (s *Stream) GetLastMove() string {
	return <-s.lastMove
}

func (s *Stream) GetLastAccuracy() float32 {
	return <-s.lastAccuracy
}

func UpdateCommandStream(start bool, accountName string,
	username1 string, username2 string, username3 string,
	sessionTimestamp uint64) {

	if start == true {
		GetStreamBuffer().Clear()
	}

	GetPositionStream().UpdateCommandStream(start)

	command1 := StreamCommand{start, username1, accountName, sessionTimestamp}
	command2 := StreamCommand{start, username2, accountName, sessionTimestamp}
	command3 := StreamCommand{start, username3, accountName, sessionTimestamp}
	GetStream(8881).commandStream <- command1
	GetStream(8882).commandStream <- command2
	GetStream(8883).commandStream <- command3
}

//checkCommandStream updates the stream username and sessionTimestamp via commandStream
func (s *Stream) checkCommandStream() bool {
	if len(s.commandStream) == 0 {
		return false
	}

	streamCommand := <-s.commandStream
	changed := s.start != streamCommand.start ||
		s.username != streamCommand.username ||
		s.sessionTimestamp != streamCommand.sessionTimestamp
	s.start = streamCommand.start
	s.username = streamCommand.username
	s.sessionTimestamp = streamCommand.sessionTimestamp

	return changed
}

//func (s *Stream) ClientListen() {
//	connPort := fmt.Sprint(s.port)
//	fmt.Println("Listening to Comms-Ultra96 via " + connType + " on " + connHost + ":" + connPort)
//	l, err := net.Listen(connType, connHost+":"+connPort)
//	if err != nil {
//		log.Fatal("Error listening:", err.Error())
//	}
//	defer l.Close()
//
//	for {
//		c, err := l.Accept()
//		if err != nil {
//			fmt.Println("Error connecting:", err.Error())
//			return
//		}
//		fmt.Println("Client connected.")
//
//		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")
//
//		go s.handleRequest(c)
//	}
//}

func (s *Stream) ClientListen() {
	connPort := fmt.Sprint(s.port)
	fmt.Println("Listening to Comms-Ultra96 via " + connType + " on " + connHost + ":" + connPort)
	c, err := net.Dial(connType, connHost+":"+connPort)
	for err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	defer c.Close()

	fmt.Println("Successfully Connected to Ultra96 on port:" + connPort)
	s.handleRequest(c)
}

func (s *Stream) ReadAlert() session.Alert{
	return <-s.alert
}

func (s *Stream) handleRequest(conn net.Conn) {
	br := bufio.NewReader(conn)
	moveNum := 0

	for {
		packetData, err := br.ReadBytes('\x7F')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}

		packet := &session.Packet{}
		err = proto.Unmarshal(packetData, packet)
		if err != nil {
			//if not session.Packet, must be alert packet
			fmt.Println("Unmarshall Error", err.Error())
			continue
		}

		if packet.DanceMove == "\x7F" {
			alert := &session.Alert{}
			if err2 := proto.Unmarshal(packetData, alert); err2 != nil {
				fmt.Println("Unmarshall Error", err2.Error())
				continue
			}

			s.alert <- *alert
			continue
		}

		if s.checkCommandStream() {
			moveNum = 0
		}

		if s.start {
			packet = confidenceLevelAdjustment(packet)

			if packet.Accuracy > -4000 {

				portMapMutex.Lock()
				GetStreamBuffer().PortMap[s.port] = append(GetStreamBuffer().PortMap[s.port], *packet)
				portMapMutex.Unlock()
			}

			//uncomment for default
			//go GetStreamBuffer().UpdateGroupSyncDelay()

			//comment these for default --> calculate true group sync delay and prompt dashboard
			go GetStreamBuffer().UpdateStartGroupSyncDelay(*packet, s.port)
			//------------------

			s.packetStream <- *packet
			moveNum += 1
		}
	}
}

func confidenceLevelAdjustment(packet *session.Packet) *session.Packet {
	if packet.Accuracy < 0 {
		return packet
	}
	if packet.Accuracy >= lowerBound3star {
		(*packet).Accuracy = 3
	} else if packet.Accuracy >= lowerBound2star {
		(*packet).Accuracy = 2
	} else {
		(*packet).Accuracy = 1
	}

	return packet
}