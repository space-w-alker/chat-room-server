package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/space-w-alker/chat-room-server/auth"
	"github.com/space-w-alker/chat-room-server/database"
	"github.com/space-w-alker/chat-room-server/model/user"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	database.InitDB()
}

func main() {
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"response": "done"})
	})
	user.RegisterHandlers(r)
	auth.RegisterHandlers(r)
	err := r.Run()
	log.Fatal(err)
}
