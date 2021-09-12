package api

import (
	"dashboard-server/comms"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

//http://arlimus.github.io/articles/gin.and.gorilla/

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	var packet comms.Packet

	for {
		packet = comms.ReadStream()
		//t, msg, err := conn.ReadMessage()
		//if err != nil {
		//	break
		//}
		conn.WriteMessage(websocket.TextMessage, []byte(packet.String()))
	}
}