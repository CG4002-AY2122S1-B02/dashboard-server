package api

import (
	"dashboard-server/internal/session/vo"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetDanceDuration(c *gin.Context) {
	var req vo.GetUserDanceReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{"message": "invalid request", "success": false})
		return
	}

	resp, err := vo.GetDanceDuration(req)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error(), "success": false})
		return
	}

	c.JSON(200, resp)
}

func (s *Server) GetDancePerformance(c *gin.Context) {
	var req vo.GetUserDanceReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{"message": "invalid request", "success": false})
		return
	}

	resp, err := vo.GetDancePerformance(req)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error(), "success": false})
		return
	}

	c.JSON(200, resp)
}

func (s *Server) GetDanceBuddies(c *gin.Context) {
	var req vo.GetUserDanceReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{"message": "invalid request", "success": false})
		return
	}

	resp, err := vo.GetDanceBuddies(req)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error(), "success": false})
		return
	}

	c.JSON(200, resp)
}
