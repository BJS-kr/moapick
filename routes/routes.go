package routes

import (
	"moapick/article"
	"moapick/tag"
	"moapick/user"

	"github.com/gofiber/fiber/v2"
)

func SetRouters(r *fiber.App)  {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
	
	user.UserRouter(r)
	article.ArticleRouter(r)
	tag.TagRouter(r)
}
