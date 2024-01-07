package article

import (
	"moapick/db"
	"moapick/db/models"

	"gorm.io/gorm/clause"
)

func SaveArticle(articleEntity *models.Article) error {
	result := db.Client.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "email"}, {Name: "title"}},
		DoUpdates: clause.AssignmentColumns([]string{"article_link", "og_image_link", "updated_at", "deleted_at"}),
	}).Create(articleEntity)

	return result.Error
}

func FindArticlesByEmail(email string) ([]models.Article, error) {
	var articles []models.Article

	result := db.Client.Find(&articles, "email = ? AND deleted_at = ?", email, nil)

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

func DeleteArticlesByEmail(email string) error {
	result := db.Client.Delete(&models.Article{}, "email = '?'", email)
	
	return result.Error
}

func UpdateArticleTitleById(articleId uint, title string) error {
	result := db.Client.Model(&models.Article{}).Where("id = ?", articleId).Update("title", title)

	return result.Error
}