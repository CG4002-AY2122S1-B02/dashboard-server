package api

import "github.com/gin-gonic/gin"

func (s *Server) GetSession(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get session",
	})
}
