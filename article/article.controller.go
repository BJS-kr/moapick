package article

import (
	"log"
	"moapick/db/models"
	"moapick/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/opengraph"
)

type SaveArticleBody struct {
	Link string `json:"link"`
	Title string `json:"title"`
}

type UpdateArticleTitleBody struct {
	Title string `json:"title"`
}

func ArticleController(r *fiber.App) {
	a := r.Group("/article")
	a.Use(middleware.JwtMiddleware())

	a.Post("/", func(c *fiber.Ctx) error {
		email, ok := c.Locals("email").(string)

		if !ok {
			log.Println("failed to assert email as string")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get email"})
		}

		article := new(SaveArticleBody)

		if err := c.BodyParser(article); err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unexpected request body"})
		
		}

		isValidUrl := IsValidURL(article.Link)

		if !isValidUrl {
			log.Println("invalid url")
			return c.Status(fiber.StatusBadRequest).JSON( fiber.Map{"error": "invalid url"})
		
		}

		articleEntity := models.Article{
		Email:       email,
		Title:       article.Title,
		ArticleLink: article.Link,
		}

		og, err := opengraph.Fetch(article.Link)

		if err == nil {
			if len(og.Image) > 0 {
				articleEntity.OgImageLink = og.Image[0].URL
			}
		}

		saveErr := SaveArticle(&articleEntity)

		if saveErr != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to save article"})
			
		}

		return c.JSON( articleEntity)
	})

	a.Get("/all", func(c *fiber.Ctx) error {
		email, ok := c.Locals("email").(string)

		if !ok {
			log.Println("failed to assert email as string")
			return c.Status(fiber.StatusInternalServerError).JSON( fiber.Map{"error": "failed to get email"})
		}

		articles, err := FindArticlesByEmail(email)

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON( fiber.Map{"error": "failed to get articles"})
		
		}

		return c.JSON( articles)
	})

	a.Get("/:articleId", func(c *fiber.Ctx)error {
		
		articleId, err := strconv.Atoi(c.Params("articleId"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON( fiber.Map{"error":"articleId must be integer"})
			
		}

		article, err := FindArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON( fiber.Map{"error": "failed to get article"})
		
		}

		return c.JSON(article)
	})

	a.Delete("/:articleId", func(c *fiber.Ctx)error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "articleId must be integer"})
			
		}

		err = DeleteArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON( fiber.Map{"error": "failed to delete article"})
			
		}

		return c.SendStatus(fiber.StatusOK)
	})

	a.Patch("/title/:articleId", func (c *fiber.Ctx) error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"articleId must be integer"})
			
		}

		updateArticleTitleBody := new(UpdateArticleTitleBody)

		if err := c.BodyParser(updateArticleTitleBody); err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON( fiber.Map{"error": "unexpected request body"})
			
		}

		updateErr := UpdateArticleTitleById(uint(articleId), updateArticleTitleBody.Title)

		if updateErr != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to update article title"})
		
		}

		return c.SendStatus(fiber.StatusOK)
	})
}
