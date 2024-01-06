package routes

import (
	"moapick/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Default에 panic recovery가 포함되어 있음
	r := gin.Default()

	user.UserController(r)

	return r
}
