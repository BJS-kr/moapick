package tag

import (
	"log"
	"moapick/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TagBody struct {
	Title string `json:"title"`
}

type ArticleIdAndTagId struct {
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

		attachBody := new(ArticleIdAndTagId)

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
		detachBody := new(ArticleIdAndTagId)

		if err := c.BodyParser(detachBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unexpected request body"})
		}

		if detachErr := DetachTagFromArticle(detachBody); detachErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to detach a tag"})
		}

		return c.SendStatus(fiber.StatusOK)
	})

	t.Delete("/:tagId", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert userId as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user id"})
		}

		tagId, err := strconv.Atoi(c.Params("tagId"))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "failed to get tag id"})
		}

		belongsToUser, err := IsTagBelongsToUser(userId, uint(tagId))

		if !belongsToUser {
			if err != nil {
				log.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete tag"})
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "tag does not belongs to user"})
			}
		}

		deleteErr := DeleteTagAndItsAssociations(uint(tagId))

		if deleteErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete tag"})
		}

		return c.SendStatus(fiber.StatusOK)
	})

	t.Get("/articles/:tagId", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert userId as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user id"})
		}
		
		tagId, err := strconv.Atoi(c.Params("tagId"))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "failed to get tag id"})
		}

		belongsToUser, err := IsTagBelongsToUser(userId, uint(tagId))

		if !belongsToUser {
			if err != nil {
				log.Println(err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get articles by tag"})
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "tag does not belongs to user"})
			}
		}
	
		articlesByTag, err := GetArticlesByTagId(uint(tagId))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"failed to find articles by tag"})
		}

		return c.Status(fiber.StatusOK).JSON(articlesByTag)
	})
}