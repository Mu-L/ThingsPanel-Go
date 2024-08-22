package sse

import (
	"fmt"
	"net/http"
	"project/api"
	"project/global"
	"project/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SSEApi struct{}

// api/v1/events

func (s *SSEApi) GetSystemEvents(c *gin.Context) {
	userClaims, ok := c.MustGet("claims").(*utils.UserClaims)
	if !ok {
		api.ErrorHandler(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	logrus.WithFields(logrus.Fields{
		"tenantID":  userClaims.TenantID,
		"userEmail": userClaims.Email,
	}).Info("User connected to SSE")

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	clientID := global.TPSSEManager.AddClient(userClaims.TenantID, userClaims.ID, c.Writer)
	defer global.TPSSEManager.RemoveClient(userClaims.TenantID, clientID)
	// 发送初始成功消息
	c.SSEvent("message", "Connected to system events")
	c.Writer.Flush()

	// 等待客户端断开连接
	<-c.Request.Context().Done()

	logrus.WithFields(logrus.Fields{
		"tenantID":  userClaims.TenantID,
		"userEmail": userClaims.Email,
	}).Info("User disconnected from SSE")
}