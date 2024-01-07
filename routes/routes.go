package routes

import (
	"moapick/article"
	"moapick/user"

	"github.com/gofiber/fiber/v2"
)

func SetRouters(r *fiber.App)  {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	
	user.UserController(r)
	article.ArticleController(r)
}
