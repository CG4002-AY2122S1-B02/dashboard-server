package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

//https://programmer.group/golang-gin-framework-with-websocket.html

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	s := &Server{}

	//Route not found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	//initialise comms server here
	//Todo

	//initialise and begin real time streaming webhooks here
	//Todo
	streamAPI := r.Group("/apis/stream")
	streamAPI.GET("/ws", func(c *gin.Context) {
		websocketHandler(c.Writer, c.Request)
	})

	//user apis
	//Todo
	//login user via sessions forever until they decide to logout
	userAPI := r.Group("/apis/user")
	userAPI.POST("/login", s.Login)

	//batch session data CRUD apis
	//Todo
	sessionAPI := r.Group("/apis/session")
	sessionAPI.GET("/get", s.GetSession)

	return &Router{r}
}

func (r *Router) Run() {
	port := 8081

	fmt.Printf("running Dashboard Server API on localhost: %v\n", port)
	err := r.Engine.Run(fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("router_start_failed:", err.Error())
	}
}

// Server interface holds all our api and some configuration
type Server struct {}