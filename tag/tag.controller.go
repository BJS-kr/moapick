package tag

import (
	"log"
	"moapick/common"
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

type TagController struct {
	TagRepository
}

// ShowAccount godoc
//
//	@Summary		create tag
//	@Description	user의 custom tag를 생성합니다.
//	@Tags			tag
//	@Accept			json
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			request			body	TagBody	true	"information to create tag"
//	@Success		201
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/tag [post]
func (tc TagController) CreateTag(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)

	if !ok {
		log.Println("failed to assert user as uint")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get userId"})
	}

	tagBody := new(TagBody)

	if err := c.BodyParser(tagBody); err != nil {
		log.Panicln(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "unexpected request body"})
	}

	err := tc.TagRepository.CreateTag(tagBody.Title, userId)

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to save tag"})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// ShowAccount godoc
//
//	@Summary	user의 모든 tag를 반환합니다.
//	@Tags		tag
//	@Param		Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success	200				{array}		models.Tag
//	@Failure	400				{object}	common.ErrorMessage
//	@Failure	500				{object}	common.ErrorMessage
//	@Router		/tag/all [get]
func (tc TagController) GetAllTagsOfUser(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)

	if !ok {
		log.Println("failed to assert userId as uint")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get user id"})
	}

	tags, err := tc.TagRepository.GetAllTagsOfUser(userId)

	if err != nil {
		log.Println(err.Error())
		c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get user tags"})
	}

	return c.JSON(tags)
}

// ShowAccount godoc
//
//	@Summary		attach tag
//	@Description	user의 custom tag를 article에 붙입니다.
//	@Tags			tag
//	@Accept			json
//	@Param			Authorization	header	string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			request			body	ArticleIdAndTagId	true	"information to attach tag"
//	@Success		200
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/tag/attach [patch]
func (tc TagController) AttachTagToArticle(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)

	if !ok {
		log.Println("failed to assert user as uint")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get userId"})
	}

	attachBody := new(ArticleIdAndTagId)

	if err := c.BodyParser(attachBody); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "unexpected request body"})
	}

	belongsToUser, err := tc.TagRepository.IsTagBelongsToUser(userId, attachBody.TagId)

	if !belongsToUser {
		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to attach tag"})
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorMessage{Error: "tag does not belongs to user"})
		}
	}

	if attachErr := tc.TagRepository.AttachTagToArticle(attachBody); attachErr != nil {
		log.Println(attachErr.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to attach tag"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ShowAccount godoc
//
//	@Summary		detach tag
//	@Description	article로부터 tag를 떼어냅니다.
//	@Tags			tag
//	@Accept			json
//	@Param			Authorization	header	string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			request			body	ArticleIdAndTagId	true	"information to detach tag"
//	@Success		200
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/tag/detach [patch]
func (tc TagController) DetachTagFromArticle(c *fiber.Ctx) error {
	detachBody := new(ArticleIdAndTagId)

	if err := c.BodyParser(detachBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorMessage{Error: "unexpected request body"})
	}

	if detachErr := tc.TagRepository.DetachTagFromArticle(detachBody); detachErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to detach a tag"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ShowAccount godoc
//
//	@Summary		delete tag
//	@Description	user의 custom tag를 삭합니다.
//	@Tags			tag
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			tagId			path	integer	true	"tag id to be deleted"
//	@Success		200
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/tag/{tagId} [delete]
func (tc TagController) DeleteTagById(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)

	if !ok {
		log.Println("failed to assert userId as uint")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get user id"})
	}

	tagId, err := strconv.Atoi(c.Params("tagId"))

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorMessage{Error: "failed to get tag id"})
	}

	belongsToUser, err := tc.TagRepository.IsTagBelongsToUser(userId, uint(tagId))

	if !belongsToUser {
		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to delete tag"})
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorMessage{Error: "tag does not belongs to user"})
		}
	}

	deleteErr := tc.TagRepository.DeleteTagAndItsAssociations(uint(tagId))

	if deleteErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to delete tag"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ShowAccount godoc
//
//	@Summary		특정 tag에 해당하는 articles를 반환합니다.
//	@Description	tag id에 해당하는 모든 articles를 반환합니다.
//	@Tags			tag
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			tagId			path		integer	true	"tag id to get articles"
//	@Success		200				{array}		models.Article
//	@Failure		400				{object}	common.ErrorMessage
//	@Failure		500				{object}	common.ErrorMessage
//	@Router			/tag/articles{tagId} [get]
func (tc TagController) GetArticlesByTagId(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(uint)

	if !ok {
		log.Println("failed to assert userId as uint")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get user id"})
	}

	tagId, err := strconv.Atoi(c.Params("tagId"))

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorMessage{Error: "failed to get tag id"})
	}

	belongsToUser, err := tc.TagRepository.IsTagBelongsToUser(userId, uint(tagId))

	if !belongsToUser {
		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to get articles by tag"})
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorMessage{Error: "tag does not belongs to user"})
		}
	}

	articlesByTag, err := tc.TagRepository.GetArticlesByTagId(uint(tagId))

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorMessage{Error: "failed to find articles by tag"})
	}

	return c.Status(fiber.StatusOK).JSON(articlesByTag)
}
