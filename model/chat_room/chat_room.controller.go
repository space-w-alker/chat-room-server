package chat_room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/space-w-alker/chat-room-server/database"
	"github.com/space-w-alker/chat-room-server/model/generic"
)

func RegisterHandlers(r *gin.Engine) {
	group := r.Group("/chat_room")
	group.POST("", CreateHandler)
	group.GET("", ReadHandler)
	group.GET("/:id", ReadOneHandler)
	group.PATCH("/:id", UpdateHandler)
	group.DELETE("/:id", DeleteHandler)
}

func CreateHandler(c *gin.Context) {
	var dto ChatRoom
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := Create(database.DB, &dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func ReadOneHandler(c *gin.Context) {
	id := c.Param("id")
	u, err := GetById(database.DB, &id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, u)
	}
}

func ReadHandler(c *gin.Context) {
	var dto ChatRoom
	var opts generic.PaginationArgs
	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBind(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, meta, err := GetWhere(database.DB, &dto, &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"items": items, "meta": meta})
	}
}

func UpdateHandler(c *gin.Context) {
	id := c.Param("id")
	var dto ChatRoom
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := Update(database.DB, &id, &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

func DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := Delete(database.DB, &id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
