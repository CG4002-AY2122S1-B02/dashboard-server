package main

import (
	"dashboard-server/api"
	"dashboard-server/comms"
	"fmt"
)

func main() {
	fmt.Println("Running Dashboard Server")
	comms.InitComm()
	go comms.ClientListen()
	//api.Run()
	api.NewRouter().Run()
}
