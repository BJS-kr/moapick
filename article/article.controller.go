package article

import (
	"log"
	"moapick/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
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

		savedArticle, err := SaveArticle(email, article)

		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to save article"})
			return
		}

		c.JSON(http.StatusCreated, *savedArticle)
	})

	a.GET("/all/:userEmail", func(ctx *gin.Context) {})
	a.GET("/:userEmail/:articleId", func(ctx *gin.Context) {})
}
