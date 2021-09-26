package api

import (
	"dashboard-server/comms"
	"dashboard-server/dbutils"
	sessionPo "dashboard-server/internal/session/po"
	streamPo "dashboard-server/internal/stream/po"
	userPo "dashboard-server/internal/user/po"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

//https://programmer.group/golang-gin-framework-with-websocket.html

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.New()
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
	streamAPI.GET("/ws/:port/*attribute", func(c *gin.Context) {
		//move this to stream file
		portStr := c.Param("port")
		attribute := c.Param("attribute")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			c.String(404, "tcp port does not exist err: "+err.Error())
			return
		}
		if comms.GetStream(port) == nil {
			c.String(404, "stream does not exist")
			return
		}
		websocketHandler(c.Writer, c.Request, port, attribute)
	})
	streamAPI.GET("/position", func(c *gin.Context) {
		websocketPositionData(c.Writer, c.Request)
	})
	streamAPI.POST("/command", s.PostStreamCommand)
	//user apis
	//Todo
	//login user via sessions forever until they decide to logout
	AccountAPI := r.Group("/apis/account")
	AccountAPI.POST("/login", s.Login)
	AccountAPI.POST("/logout", s.Logout)
	UserAPI := AccountAPI.Group("/users")
	UserAPI.POST("/register", s.RegisterUsers)

	//batch session data CRUD apis
	//Todo
	sessionAPI := r.Group("/apis/session")
	sessionAPI.GET("/get", s.GetSession)
	sessionAPI.GET("/current", s.GetCurrentSession)
	sessionAPI.POST("/upload", s.UploadSession)

	danceAPI := r.Group("/apis/dance")
	danceAPI.GET("/duration", s.GetDanceDuration)
	danceAPI.GET("/buddies", s.GetDanceBuddies)
	danceAPI.GET("/performance", s.GetDancePerformance)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}

	r.Use(cors.New(config))
	return &Router{r}
}

func (r *Router) Run() {
	err := dbutils.GetDB().AutoMigrate(
		&streamPo.SensorData{},
		&streamPo.SyncDelay{},
		&sessionPo.Session{},
		&sessionPo.UserSession{},
		&userPo.User{},
		&userPo.Account{},
	)
	if err != nil {
		log.Fatal("db auto migrate failed:", err.Error())
	}

	comms.GetPositionStream()
	comms.NewStream(8881)
	comms.NewStream(8882)
	comms.NewStream(8883)

	port := 8081

	fmt.Printf("running Dashboard Server API on localhost: %v\n", port)
	err = r.Engine.Run(fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("router_start_failed:", err.Error())
	}
}

// Server interface holds all our api and some configuration
type Server struct{}
