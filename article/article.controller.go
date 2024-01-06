package article

import (
	"log"
	"moapick/db/models"
	"moapick/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/opengraph"
)

type ArticleBody struct {
	Link string `json:"link"`
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

		article := ArticleBody{}

		if err := c.ShouldBind(&article); err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "unexpected request type"})
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
		ArticleLink:article.Link,
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

	a.GET("/all/:userEmail", func(ctx *gin.Context) {})
	a.GET("/:userEmail/:articleId", func(ctx *gin.Context) {})
}
