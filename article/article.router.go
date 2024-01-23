package article

import (
	"moapick/db"
	"moapick/middleware"

	"github.com/gofiber/fiber/v2"
)

type SaveArticleBody struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

type UpdateArticleTitleBody struct {
	Title string `json:"title"`
}

func ArticleRouter(r *fiber.App) {
	articleRepository := ArticleRepository{Client: db.Client}
	articleService := ArticleService{}
	articleController := ArticleController{ArticleRepository: articleRepository, ArticleService: articleService}

	a := r.Group("/article")
	a.Use(middleware.JwtMiddleware())

	a.Post("/", articleController.SaveArticle)
	a.Get("/all", articleController.GetAllArticlesOfUser)
	a.Get("/:articleId", articleController.GetArticleById)
	a.Delete("/all", articleController.DeleteArticlesByUserId)
	a.Delete("/:articleId", articleController.DeleteArticleById)
	a.Patch("/title/:articleId", articleController.UpdateArticleTitleById)
}
