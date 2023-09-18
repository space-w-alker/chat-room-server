package user

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateHandler(c *gin.Context)  {
	var dto CreateUserDTO
	c.BindJSON(&dto)
	Create(&dto)
}

func ReadOneHandler(c *gin.Context)  {
	pathSplit := strings.Split(c.FullPath(), "/")
	id := pathSplit[len(pathSplit)-1]
	GetById(&id)
}
func ReadHandler(c *gin.Context)  {
	var dto GetOrUpdateUserDTO
	c.BindJSON(&dto)
	GetWhere(&dto)
}

func UpdateHandler(c *gin.Context)  {
	pathSplit := strings.Split(c.FullPath(), "/")
	id := pathSplit[len(pathSplit)-1]
	var dto GetOrUpdateUserDTO
	c.BindJSON(&dto)
	Update(&id, &dto)
}

func DeleteHandler(c *gin.Context)  {
	pathSplit := strings.Split(c.FullPath(), "/")
	id := pathSplit[len(pathSplit)-1]
	Delete(&id)
}