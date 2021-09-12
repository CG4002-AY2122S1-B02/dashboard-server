package api

import "github.com/gin-gonic/gin"

func (s *Server) Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "login",
	})
}
