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

	lowerBound3star = 0.65
	lowerBound2star = 0.55
)

var (
	streamMap        map[int]*Stream
	streamBuffer     *StreamBuffer
	streamBufferOnce sync.Once
)

type StreamCommand struct {
	start            bool
	username         string
	accountName      string
	sessionTimestamp uint64
}

type StreamBuffer struct {
	PortMap  map[int][]session.Packet
	Position       []session.Position
	groupSyncDelay chan uint64
	pointer        int
	totalSyncDelay uint64
}

func (sb *StreamBuffer) Clear() {
	sb.PortMap = make(map[int][]session.Packet)
	sb.Position = make([]session.Position, 0)
	sb.groupSyncDelay = make(chan uint64, bufferLength)
	sb.pointer = 0
	sb.totalSyncDelay = 0
}

func (sb *StreamBuffer) ReadGroupSyncDelay() uint64 {
	return <- sb.groupSyncDelay
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
	sb.totalSyncDelay += syncDelay

	sb.pointer += 1
}

func (sb *StreamBuffer) GetAvgSyncDelay() uint64 {
	if sb.pointer == 0 {
		return 0
	}

	return uint64(math.Round(float64(sb.totalSyncDelay) / float64(sb.pointer)))
}

func GetStreamBuffer() *StreamBuffer {
	streamBufferOnce.Do(func() {
		streamBuffer = &StreamBuffer{
			make(map[int][]session.Packet),
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
			fmt.Println("Unmarshall Error", err.Error())
		}

		if s.checkCommandStream() {
			moveNum = 0
		}

		if s.start {
			packet = confidenceLevelAdjustment(packet)

			if packet.DanceMove != "START" && packet.Accuracy > -4000 {
				GetStreamBuffer().PortMap[s.port] = append(GetStreamBuffer().PortMap[s.port], *packet)
			}

			go GetStreamBuffer().UpdateGroupSyncDelay()
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