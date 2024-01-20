package article

import (
	"moapick/db"
	"moapick/db/models"

	"gorm.io/gorm/clause"
)

func SaveArticle(userId uint, saveArticleBody *SaveArticleBody, ogImageLink string) error {
	articleEntity := models.Article{
		UserId:      userId,
		Title:       saveArticleBody.Title,
		ArticleLink: saveArticleBody.Link,
		OgImageLink: ogImageLink,
	}
	// 아티클 저장 기록 삭제시 soft delete이기 때문에 upsert로 처리해야한다.
	// articles table은 email, title이 unique index이기 때문이다.
	result := db.Client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "title"}},
		DoUpdates: clause.AssignmentColumns([]string{"article_link", "og_image_link", "updated_at", "deleted_at"}),
	}).Create(&articleEntity)

	return result.Error
}

func FindArticlesByUserId(userId uint) ([]models.Article, error) {
	var articles []models.Article
	result := db.Client.Where("user_id = ?", userId).Preload("Tags").Find(&articles)

	return articles, result.Error
}

func FindArticleById(articleId uint) (models.Article, error) {
	var article models.Article
	result := db.Client.Where("id = ?", articleId).Preload("Tags").First(&article)

	return article, result.Error
}

func DeleteArticleById(articleId uint) error {
	result := db.Client.Where("id = ?", articleId).Delete(&models.Article{})

	return result.Error
}

func DeleteArticlesByUserId(userId uint) error {
	result := db.Client.Where("user_id = ?", userId).Delete(&models.Article{})

	return result.Error
}

func UpdateArticleTitleById(articleId uint, title string) error {
	result := db.Client.Model(&models.Article{}).Where("id = ?", articleId).Update("title", title)

	return result.Error
}
