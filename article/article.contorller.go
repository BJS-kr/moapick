package article

import (
	"fmt"
	"log"
	"moapick/common"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/otiai10/opengraph"
)

type ArticleController struct {
	ArticleRepository
	ArticleService
}

// ShowAccount godoc
//
//	@Summary		save article
//	@Description	user에 속한 article을 저장합니다. OG image를 탐색하고 없을 시 빈 값이 저장됩니다.
//	@Tags			article
//	@Accept			json
//	@Param			Authorization	header	string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			request			body	SaveArticleBody	true	"information to save article"
//	@Success		201
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/article [post]
func (ac ArticleController) SaveArticle(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)

	if !ok {
		log.Println("failed to assert user as uint")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get userId"})
	}

	articleBody := new(SaveArticleBody)

	if err := c.BodyParser(articleBody); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "unexpected request body"})
	}

	isValidUrl := ac.ArticleService.IsValidURL(articleBody.Link)

	if !isValidUrl {
		log.Println("invalid url")
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "invalid url"})
	}

	ogImageLink := ""

	og, err := opengraph.Fetch(articleBody.Link)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err == nil {
		fmt.Println(og.Image)
		if len(og.Image) > 0 {
			ogImageLink = og.Image[0].URL
		}
	}

	saveErr := ac.ArticleRepository.SaveArticle(userId, articleBody, ogImageLink)

	if saveErr != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "failed to save article"})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// ShowAccount godoc
//
//	@Summary		get all articles
//	@Description	user에 속한 모든 articles를 반환합니다.
//	@Tags			article
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{array}		models.Article
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/article/all [get]
func (ac ArticleController) GetAllArticlesOfUser(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)
	if !ok {
		log.Println("failed to assert email as string")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get email"})
	}

	articles, err := ac.ArticleRepository.FindArticlesByUserId(userId)

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get articles"})
	}

	return c.JSON(articles)
}

// ShowAccount godoc
//
//	@Summary		get article by article id
//	@Description	article id에 해당하는 article을 반환합니다.
//	@Tags			article
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			articleId		path		integer	true	"article id"
//	@Success		200				{object}	models.Article
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/article/{articleId} [get]
func (ac ArticleController) GetArticleById(c *fiber.Ctx) error {
	articleId, err := strconv.Atoi(c.Params("articleId"))

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "articleId must be integer"})
	}

	article, err := ac.ArticleRepository.FindArticleById(uint(articleId))

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get article"})
	}

	return c.JSON(article)
}

// ShowAccount godoc
//
//	@Summary		delete all article
//	@Description	user에 속한 모든 articles를 지웁니다.
//	@Tags			article
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200
//	@Failure		400	{object}	common.ErrorMessage
//	@Failure		500	{object}	common.ErrorMessage
//	@Router			/article/all [delete]
func (ac ArticleController) DeleteArticlesByUserId(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)

	if !ok {
		log.Println("failed to assert userId as uint")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get userId"})
	}

	err := ac.ArticleRepository.DeleteArticlesByUserId(userId)

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to delete articles"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ShowAccount godoc
//
//	@Summary		delete article by id
//	@Description	user에 속한 article을 지웁니다.
//	@Tags			article
//	@Param			articleId		path	integer	true	"article id"
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200
//	@Failure		400	{object}	common.ErrorMessage
//	@Failure		500	{object}	common.ErrorMessage
//	@Router			/article/{articleId} [delete]
func (ac ArticleController) DeleteArticleById(c *fiber.Ctx) error {
	articleId, err := strconv.Atoi(c.Params("articleId"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "articleId must be integer"})
	}

	err = ac.ArticleRepository.DeleteArticleById(uint(articleId))

	if err != nil {
		log.Println(err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to delete article"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ShowAccount godoc
//
//	@Summary		update article title
//	@Description	article의 title을 업데이트합니다.
//	@Tags			article
//	@Accept			json
//	@Param			Authorization	header	string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			request			body	UpdateArticleTitleBody	true	"title to update"
// 	@Param          articleId       path    integer                 true    "article id"
//	@Success		200
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/article/title/{articleId} [patch]
func (ac ArticleController) UpdateArticleTitleById(c *fiber.Ctx) error {
	articleId, err := strconv.Atoi(c.Params("articleId"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "articleId must be integer"})
	}

	updateArticleTitleBody := new(UpdateArticleTitleBody)

	if err := c.BodyParser(updateArticleTitleBody); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "unexpected request body"})
	}

	updateErr := ac.ArticleRepository.UpdateArticleTitleById(uint(articleId), updateArticleTitleBody.Title)

	if updateErr != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "failed to update article title"})
	}

	return c.SendStatus(fiber.StatusOK)
}
