package handler

import (
	"github.com/gin-gonic/gin"
)

type RoomHandler interface {
	CreateRoom(c *gin.Context)
	FindRoom(c *gin.Context)
	SearchRoom(c *gin.Context)
	UpdateRoom(c *gin.Context)
	DeleteRoom(c *gin.Context)
}
