package routes

import (
	"moapick/article"
	"moapick/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(r *fiber.App)  {

	user.UserController(r)
	article.ArticleController(r)
}
