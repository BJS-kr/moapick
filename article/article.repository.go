package article

import (
	"moapick/db"
	"moapick/db/models"
)

func SaveArticle(articleEntity *models.Article) error {
	result := db.Client.Create(articleEntity)

	return result.Error
}

func FindArticlesByEmail(email string) ([]models.Article, error) {
	var articles []models.Article

	result := db.Client.Find(&articles, "email = ?", email)

	return articles, result.Error
}


func FindArticleById(articleId uint) (models.Article, error) {
	var article models.Article

	result := db.Client.First(&article, "id = ?", articleId)

	return article, result.Error
}

func DeleteArticleById(articleId uint) error {
	result := db.Client.Delete(&models.Article{}, "id = ?", articleId)

	return result.Error
}

func UpdateArticleTitleById(articleId uint, title string) error {
	result := db.Client.Model(&models.Article{}).Where("id = ?", articleId).Update("title", title)

	return result.Error
}