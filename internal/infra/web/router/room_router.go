package router

import (
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"

	"github.com/gin-gonic/gin"
)

func RoomRouter(
	r *gin.RouterGroup,
	roomHandler handler.RoomHandler,
) {
	rooms := r.Group("/rooms")
	{
		rooms.POST("", roomHandler.CreateRoom)
		rooms.GET("", roomHandler.SearchRoom)
		rooms.GET(":id", roomHandler.FindRoom)
		rooms.PUT(":id", roomHandler.UpdateRoom)
		rooms.DELETE(":id", roomHandler.DeleteRoom)
		rooms.POST(":id/send", roomHandler.SendMessage)
	}
}
