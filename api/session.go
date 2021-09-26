package api

import (
	"dashboard-server/comms"
	"dashboard-server/internal/session/vo"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetSession(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "get session",
	})
}

func (s *Server) GetCurrentSession(c *gin.Context) {
	c.JSON(200,
		vo.GetCurrentSessionResp{
			User1:    comms.GetStreamBuffer().PortMap[8881],
			User2:    comms.GetStreamBuffer().PortMap[8882],
			User3:    comms.GetStreamBuffer().PortMap[8883],
			Position: comms.GetStreamBuffer().Position,
		})
}

func (s *Server) UploadSession(c *gin.Context) {
	var req vo.UploadSessionReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, gin.H{"message": "invalid request", "success": false})
		return
	}

	resp, err := vo.UploadSession(req)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error(), "success": false})
		return
	}

	c.JSON(200, resp)
}
