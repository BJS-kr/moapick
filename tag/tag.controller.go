package tag

import "github.com/gofiber/fiber/v2"

func TagController(r *fiber.App) {
	t := r.Group("/tag")
	
	t.Post("/", func(c *fiber.Ctx) error {
		
	})

	t.Get("/all", func(c *fiber.Ctx) error {
		
	})

	t.Patch("/attach", func(c *fiber.Ctx) error {

	})

	t.Patch("/detach", func(c *fiber.Ctx) error {

	})

	t.Delete("/:tagId", func(c *fiber.Ctx) error {

	})
}