package article

import (
	"log"
	"moapick/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/opengraph"
)

type SaveArticleBody struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

type UpdateArticleTitleBody struct {
	Title string `json:"title"`
}

func ArticleController(r *fiber.App) {
	a := r.Group("/article")
	a.Use(middleware.JwtMiddleware())

	a.Post("/", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert user as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get userId"})
		}

		articleBody := new(SaveArticleBody)

		if err := c.BodyParser(articleBody); err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unexpected request body"})
		}

		isValidUrl := IsValidURL(articleBody.Link)

		if !isValidUrl {
			log.Println("invalid url")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid url"})
		}

		ogImageLink := ""

		og, err := opengraph.Fetch(articleBody.Link)

		if err == nil {
			if len(og.Image) > 0 {
				ogImageLink = og.Image[0].URL
			}
		}

		saveErr := SaveArticle(userId, articleBody, ogImageLink)

		if saveErr != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to save article"})
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	a.Get("/all", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)
		if !ok {
			log.Println("failed to assert email as string")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get email"})
		}

		articles, err := FindArticlesByUserId(userId)

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get articles"})
		}

		return c.JSON(articles)
	})

	a.Get("/:articleId", func(c *fiber.Ctx) error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			log.Panicln(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "articleId must be integer"})
		}

		article, err := FindArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get article"})
		}

		return c.JSON(article)
	})

	a.Delete("/all", func(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert userId as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get userId"})
		}

		err := DeleteArticlesByUserId(userId)

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete articles"})
		}

		return c.SendStatus(fiber.StatusOK)
	})

	a.Delete("/:articleId", func(c *fiber.Ctx) error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "articleId must be integer"})
		}

		err = DeleteArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete article"})
		}

		return c.SendStatus(fiber.StatusOK)
	})

	a.Patch("/title/:articleId", func(c *fiber.Ctx) error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "articleId must be integer"})
		}

		updateArticleTitleBody := new(UpdateArticleTitleBody)

		if err := c.BodyParser(updateArticleTitleBody); err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unexpected request body"})
		}

		updateErr := UpdateArticleTitleById(uint(articleId), updateArticleTitleBody.Title)

		if updateErr != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to update article title"})
		}

		return c.SendStatus(fiber.StatusOK)
	})
}
