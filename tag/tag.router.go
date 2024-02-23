package tag

import (
	"moapick/db"
	"moapick/middleware"

	"github.com/gofiber/fiber/v2"
)

func TagRouter(r *fiber.App) {
	tagRepository := TagRepository{Client: db.Client}
	tagController := TagController{TagRepository: tagRepository}

	t := r.Group("/tag")
	t.Use(middleware.JwtMiddleware())

	t.Post("/", tagController.CreateTag)
	t.Get("/all", tagController.GetAllTagsOfUser)
	t.Patch("/attach", tagController.AttachTagToArticle)
	t.Patch("/detach", tagController.DetachTagFromArticle)
	t.Patch("/:tagId", tagController.UpdateTagById)
	t.Delete("/:tagId", tagController.DeleteTagById)
	t.Get("/articles/:tagId", tagController.GetArticlesByTagId)
}
