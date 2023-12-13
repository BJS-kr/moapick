package routes

import (
	"moapick/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	user.UserController(r)

	return r
}
