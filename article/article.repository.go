package article

import (
	"moapick/db"
	"moapick/db/models"

	"gorm.io/gorm/clause"
)

func SaveArticle(articleEntity *models.Article) error {
	// 아티클 저장 기록 삭제시 soft delete이기 때문에 upsert로 처리해야한다.
	// articles table은 email, title이 unique index이기 때문이다.
	result := db.Client.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "title"}},
		DoUpdates: clause.AssignmentColumns([]string{"article_link", "og_image_link", "updated_at", "deleted_at"}),
	}).Create(articleEntity)

	return result.Error
}

func FindArticlesByUserId(userId uint) ([]models.Article, error) {
	var articles []models.Article

	result := db.Client.Where("user_id = ?", userId).Preload("Tags").First(&articles)

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

func DeleteArticlesByUserId(userId uint) error {
	result := db.Client.Delete(&models.Article{}, "id = '?'", userId)
	
	return result.Error
}

func UpdateArticleTitleById(articleId uint, title string) error {
	result := db.Client.Model(&models.Article{}).Where("id = ?", articleId).Update("title", title)

	return result.Error
}