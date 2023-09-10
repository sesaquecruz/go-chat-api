package handler

import (
	"github.com/gin-gonic/gin"
)

type MessageHandler interface {
	CreateMessage(c *gin.Context)
}
