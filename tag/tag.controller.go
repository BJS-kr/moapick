package tag

import (
	"log"
	"moapick/middleware"

	"github.com/gofiber/fiber/v2"
)

type TagBody struct {
	Title string `json:"title"`
}

type AttachBody struct {
	ArticleId uint `json:"article_id"`
	TagId     uint `json:"tag_id"`
}

func TagController(r *fiber.App) {
	t := r.Group("/tag")
	t.Use(middleware.JwtMiddleware())

	t.Post("/", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert user as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get userId"})
		}

		tagBody := new(TagBody)

		if err := c.BodyParser(tagBody); err != nil {
			log.Panicln(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unexpected request body"})
		}

		err := CreateTag(tagBody.Title, userId)

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save tag"})
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	t.Get("/all", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert userId as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user id"})
		}

		tags, err := GetAllTagsOfUser(userId)

		if err != nil {
			log.Println(err.Error())
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user tags"})
		}

		return c.JSON(tags)
	})

	t.Patch("/attach", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert user as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get userId"})
		}

		attachBody := new(AttachBody)

		if err := c.BodyParser(attachBody); err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unexpected request body"})
		}

		belongsToUser, err := IsTagBelongsToUser(userId, attachBody.TagId)

		if !belongsToUser {
			if err != nil {
				log.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to attach tag"})
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "tag does not belongs to user"})
			}
		}

		if attachErr := AttachTagToArticle(attachBody); attachErr != nil {
			log.Println(attachErr.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to attach tag"})
		}

		return c.SendStatus(fiber.StatusOK)
	})

	t.Patch("/detach", func(c *fiber.Ctx) error {
		return c.SendStatus(400)
	})

	t.Delete("/:tagId", func(c *fiber.Ctx) error {
		return c.SendStatus(400)
	})
}
