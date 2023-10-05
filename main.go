package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/space-w-alker/chat-room-server/auth"
	"github.com/space-w-alker/chat-room-server/database"
	"github.com/space-w-alker/chat-room-server/model/chat"
	"github.com/space-w-alker/chat-room-server/model/chat_room"
	"github.com/space-w-alker/chat-room-server/model/room_tag"
	"github.com/space-w-alker/chat-room-server/model/user"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".local.env")
	viper.ReadInConfig()
	database.InitDB()
}

func main() {
	r := gin.Default()
	user.RegisterHandlers(r)
	auth.RegisterHandlers(r)
	chat.RegisterHandlers(r)
	chat_room.RegisterHandlers(r)
	room_tag.RegisterHandlers(r)
	r.Static("public", "./public")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	err := r.Run()
	log.Fatal(err)
}
