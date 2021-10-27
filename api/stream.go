package api

import (
	"dashboard-server/comms"
	"dashboard-server/internal/session"
	"dashboard-server/internal/stream/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

//http://arlimus.github.io/articles/gin.and.gorilla/

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, r *http.Request, port int, attribute string) error {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return errors.Wrap(err, "websocket handler error")
	}

	for {
		stream := comms.GetStream(port)
		if stream == nil {
			return errors.New("websocket handler error: stream does not exist")
		}

		switch attribute {
		case session.All.String():
			packet := stream.ReadStream()
			if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v|%v|%v", packet.EpochMs, packet.DanceMove, packet.Accuracy))); err != nil {
				return errors.Wrap(err, "failed to read 'all' stream")
			}
		case session.DanceMove.String():
			if err := conn.WriteMessage(websocket.TextMessage, []byte(stream.GetLastMove())); err != nil {
				return errors.Wrap(err, "failed to read 'dance move' stream")
			}
		case session.Accuracy.String():
			if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(stream.GetLastAccuracy()))); err != nil {
				return errors.Wrap(err, "failed to read 'accuracy' stream")
			}
		case session.EpochMs.String():
			if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("not implemented"))); err != nil {
				return errors.Wrap(err, "failed to read 'time ms' stream")
			}
		}
	}
}

func websocketAlert(w http.ResponseWriter, r *http.Request, port int) error {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return errors.Wrap(err, "websocket handler error")
	}

	for {
		stream := comms.GetStream(port)
		alert := stream.ReadAlert()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(alert.Message)); err != nil {
			return errors.Wrap(err, "failed to read 'alert' stream")
		}
	}
}

func websocketPositionData(w http.ResponseWriter, r *http.Request) error {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return errors.Wrap(err, "websocket handler error")
	}

	for {
		position := comms.PositionDataStream.ReadStream()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(position.Position)); err != nil {
			return errors.Wrap(err, "failed to read 'position' stream")
		}
	}
}

func websocketECGData(w http.ResponseWriter, r *http.Request) error {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return errors.Wrap(err, "websocket handler error")
	}

	for {
		ecgData := comms.ECGDataStream.ReadStream()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(ecgData.Val1))); err != nil {
			return errors.Wrap(err, "failed to read 'ecg data' stream")
		}
	}
}

func websocketGroupSyncDelay(w http.ResponseWriter, r *http.Request) error {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return errors.Wrap(err, "websocket handler error")
	}

	for {
		groupSyncDelay := comms.GetStreamBuffer().ReadGroupSyncDelay()

		if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint(groupSyncDelay))); err != nil {
			return errors.Wrap(err, "failed to retrieve group sync delay")
		}
	}
}

func (s *Server) PostStreamCommand(c *gin.Context) {
	var req vo.PostStreamCommandReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{"message": "invalid request", "success": false})
		return
	}

	comms.UpdateCommandStream(req.Start, req.AccountName, req.Username1, req.Username2, req.Username3, req.SessionTimestamp)

	c.JSON(200, gin.H{
		"success": true,
		"command": req,
	})
}
