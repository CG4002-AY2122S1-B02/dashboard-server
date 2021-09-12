package api

import (
	"dashboard-server/comms"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

func Echo(ws *websocket.Conn) {

	var (
		err error
		packet comms.Packet
	)

	for {
		//msg = fmt.Sprint(counter)
		packet = comms.ReadStream()

		if err = websocket.Message.Send(ws, packet.String()); err != nil {
			fmt.Println("Can't send", err.Error())
		} else {
			fmt.Println("Sending: ", packet)
		}
	}
}

func Run() {
	const maxClients = 1
	sema := make(chan struct{}, maxClients)

	http.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		sema <- struct{}{}
		defer func() { <-sema }()
		Echo(ws)
	}))

	fmt.Println("Running websocket stream on localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}