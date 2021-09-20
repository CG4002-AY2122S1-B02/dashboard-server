package api

import (
	"dashboard-server/internal/user/vo"
	"fmt"
	"github.com/gin-gonic/gin"
)

//Pre Selectable User Icons to Chose From* Like some game avatar/ upload icon*

func (s *Server) Login(c *gin.Context) {
	var req vo.LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{"message": "invalid request", "success": false})
		return
	}

	resp, err := vo.Login(req.AccountName, req.Password)
	if err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}

func (s *Server) RegisterUsers(c *gin.Context) {
	var req vo.RegisterUsersReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{"message": "invalid request", "success": false})
		return
	}

	fmt.Println(req)

	resp, err := vo.RegisterUsersIfNotExist(req)
	if err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}

func (s *Server) Logout(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "logout",
	})
}
