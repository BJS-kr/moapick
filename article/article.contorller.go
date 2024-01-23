package article

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/opengraph"
)

type ArticleController struct {
	ArticleRepository
	ArticleService
}

func (ac ArticleController)SaveArticle(c *fiber.Ctx) error {
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

		isValidUrl := ac.ArticleService.IsValidURL(articleBody.Link)

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

		saveErr := ac.ArticleRepository.SaveArticle(userId, articleBody, ogImageLink)

		if saveErr != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to save article"})
		}

		return c.SendStatus(fiber.StatusCreated)
	}

func (ac ArticleController)GetAllArticlesOfUser(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)
		if !ok {
			log.Println("failed to assert email as string")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get email"})
		}

		articles, err := ac.ArticleRepository.FindArticlesByUserId(userId)

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get articles"})
		}

		return c.JSON(articles)
	}

func(ac ArticleController)GetArticleById(c *fiber.Ctx) error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "articleId must be integer"})
		}

		article, err := ac.ArticleRepository.FindArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get article"})
		}

		return c.JSON(article)
	}

func(ac ArticleController)DeleteArticlesByUserId(c *fiber.Ctx) error {
		userId, ok := c.Locals("userId").(uint)

		if !ok {
			log.Println("failed to assert userId as uint")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get userId"})
		}

		err := ac.ArticleRepository.DeleteArticlesByUserId(userId)

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete articles"})
		}

		return c.SendStatus(fiber.StatusOK)
	}

func(ac ArticleController)DeleteArticleById(c *fiber.Ctx) error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "articleId must be integer"})
		}

		err = ac.ArticleRepository.DeleteArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete article"})
		}

		return c.SendStatus(fiber.StatusOK)
	}

func(ac ArticleController)UpdateArticleTitleById(c *fiber.Ctx) error {
		articleId, err := strconv.Atoi(c.Params("articleId"))

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "articleId must be integer"})
		}

		updateArticleTitleBody := new(UpdateArticleTitleBody)

		if err := c.BodyParser(updateArticleTitleBody); err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unexpected request body"})
		}

		updateErr := ac.ArticleRepository.UpdateArticleTitleById(uint(articleId), updateArticleTitleBody.Title)

		if updateErr != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to update article title"})
		}

		return c.SendStatus(fiber.StatusOK)
	}