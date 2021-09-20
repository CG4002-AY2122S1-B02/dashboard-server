package main

import (
	"dashboard-server/api"
	"fmt"
)

func main() {
	fmt.Println("Running Dashboard Server")
	api.NewRouter().Run()
}
