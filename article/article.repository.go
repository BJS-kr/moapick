package article

import (
	"moapick/db/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArticleRepository struct {
	Client *gorm.DB
}

func (ar ArticleRepository) SaveArticle(userId uint, saveArticleBody *SaveArticleBody, ogImageLink string) error {
	articleEntity := models.Article{
		UserId:      userId,
		Title:       saveArticleBody.Title,
		ArticleLink: saveArticleBody.Link,
		OgImageLink: ogImageLink,
	}
	// 아티클 저장 기록 삭제시 soft delete이기 때문에 upsert로 처리해야한다.
	// articles table은 email, title이 unique index이기 때문이다.
	result := ar.Client.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "title"}},
		DoUpdates: clause.AssignmentColumns([]string{"article_link", "og_image_link", "updated_at", "deleted_at"}),
	}).Create(&articleEntity)

	return result.Error
}

func (ar ArticleRepository) FindArticlesByUserId(userId uint) ([]models.Article, error) {
	var articles []models.Article
	result := ar.Client.Where("user_id = ?", userId).Preload("Tags", func(db *gorm.DB) *gorm.DB {
		// select한다고 해서 return 값에서 select된 필드만 들어가는 것은 아니다.
		// 없는 값은 zero value가 들어간다.정말 원하는 필드만 전달하기 위해선 return dto가 필수
		return db.Select("id", "title", "created_at")
	}).Find(&articles)

	return articles, result.Error
}

func (ar ArticleRepository) FindArticleById(articleId uint) (models.Article, error) {
	var article models.Article
	result := ar.Client.Where("id = ?", articleId).Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "title", "created_at")
	}).First(&article)

	return article, result.Error
}

func (ar ArticleRepository) DeleteArticleById(articleId uint) error {
	result := ar.Client.Where("id = ?", articleId).Delete(&models.Article{})

	return result.Error
}

func (ar ArticleRepository) DeleteArticlesByUserId(userId uint) error {
	result := ar.Client.Where("user_id = ?", userId).Delete(&models.Article{})

	return result.Error
}

func (ar ArticleRepository) UpdateArticleTitleById(articleId uint, title string) error {
	result := ar.Client.Model(&models.Article{}).Where("id = ?", articleId).Update("title", title)

	return result.Error
}
