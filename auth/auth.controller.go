package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/space-w-alker/chat-room-server/model/user"
)

func RegisterHandlers(engine *gin.Engine) {
	authGroup := engine.Group("/auth")
	authGroup.POST("sign-in", SignInHandler)
	authGroup.POST("sign-up", SignUpHandler)
}

func SignUpHandler(c *gin.Context) {
	var dto user.User
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := SignUp(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusAccepted, gin.H{"token": token})
	}
}

func SignInHandler(c *gin.Context) {
	var dto SignInDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := SignIn(dto.Email, dto.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusAccepted, gin.H{"token": token})
	}
}
