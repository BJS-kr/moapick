package routes

import (
	"moapick/article"
	"moapick/user"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Default에 panic recovery가 포함되어 있음
	r := gin.Default()

	user.UserController(r)
	article.ArticleController(r)

	return r
}
