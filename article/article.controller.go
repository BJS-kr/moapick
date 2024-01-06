package article

import (
	"log"
	"moapick/db/models"
	"moapick/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/opengraph"
)

type SaveArticleBody struct {
	Link string `json:"link"`
	Title string `json:"title"`
}

type UpdateArticleTitleBody struct {
	Title string `json:"title"`
}

func ArticleController(r *gin.Engine) {
	a := r.Group("/article")
	a.Use(middleware.JwtMiddleware())

	a.POST("/", func(c *gin.Context) {
		IEmail, ok := c.Get("email")

		if !ok {
			log.Println("context does not have user email when trying to insert article link")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get email"})
		}

		email, ok := IEmail.(string)

		if !ok {
			log.Println("failed to assert email as string")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get email"})
		}

		article := SaveArticleBody{}

		if err := c.ShouldBind(&article); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "unexpected request body"})
			return
		}

		isValidUrl := IsValidURL(article.Link)

		if !isValidUrl {
			log.Println("invalid url")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
			return
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to save article"})
			return
		}

		c.JSON(http.StatusCreated, articleEntity)
	})

	a.GET("/all", func(c *gin.Context) {
		IEmail, ok := c.Get("email")

		if !ok {
			log.Println("context does not have user email when trying to insert article link")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get email"})
		}

		email, ok := IEmail.(string)

		if !ok {
			log.Println("failed to assert email as string")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get email"})
		}

		articles, err := FindArticlesByEmail(email)

		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get articles"})
			return
		}

		c.JSON(http.StatusOK, articles)
	})

	a.GET("/:articleId", func(c *gin.Context) {
		
		articleId, err := strconv.Atoi(c.Param("articleId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "articleId must be integer")
			return
		}

		article, err := FindArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get article"})
			return
		}

		c.JSON(http.StatusOK, article)
	})

	a.DELETE("/:articleId", func(c *gin.Context) {
		articleId, err := strconv.Atoi(c.Param("articleId"))

		if err != nil {
			c.JSON(http.StatusBadRequest, "articleId must be integer")
			return
		}

		err = DeleteArticleById(uint(articleId))

		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete article"})
			return
		}

		c.Status(http.StatusOK)
	})

	a.PATCH("/title/:articleId", func (c *gin.Context) {
		articleId, err := strconv.Atoi(c.Param("articleId"))

		if err != nil {
			c.JSON(http.StatusBadRequest, "articleId must be integer")
			return
		}

		updateArticleTitleBody := UpdateArticleTitleBody{}

		if err := c.ShouldBind(&updateArticleTitleBody); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "unexpected request body"})
			return
		}

		updateErr := UpdateArticleTitleById(uint(articleId), updateArticleTitleBody.Title)

		if updateErr != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update article title"})
		return
		}

		c.Status(http.StatusOK)
	})
}
