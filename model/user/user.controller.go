package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(r *gin.Engine)  {
	group := r.Group("/user")
	group.POST("/", CreateHandler)
	group.GET("/", ReadHandler)
	group.GET("/:id", ReadOneHandler)
	group.PATCH("/:id", UpdateHandler)
	group.DELETE("/:id", DeleteHandler)
}

func CreateHandler(c *gin.Context)  {
	var dto CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Create(&dto)
}

func ReadOneHandler(c *gin.Context)  {
	id := c.Param("id")
	GetById(&id)
}

func ReadHandler(c *gin.Context)  {
	var dto GetOrUpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	GetWhere(&dto)
}

func UpdateHandler(c *gin.Context)  {
	id := c.Param("id")
	var dto GetOrUpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Update(&id, &dto)
}

func DeleteHandler(c *gin.Context)  {
	id := c.Param("id")
	Delete(&id)
}