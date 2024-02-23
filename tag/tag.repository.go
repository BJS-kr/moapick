package tag

import (
	"errors"
	"moapick/db/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	Client *gorm.DB
}

func (tr TagRepository) CreateTag(title string, userId uint) error {
	tagEntity := models.Tag{
		Title:  title,
		UserID: userId,
	}

	result := tr.Client.Create(&tagEntity)

	return result.Error
}

func (tr TagRepository) GetAllTagsOfUser(userId uint) ([]models.Tag, error) {
	var tags []models.Tag
	result := tr.Client.Where("user_id = ?", userId).Find(&tags)

	return tags, result.Error
}

func (tr TagRepository) IsTagBelongsToUser(userId, tagId uint) (bool, error) {
	if result := tr.Client.Select("id").Where("id = ? AND user_id = ?", tagId, userId).First(&models.Tag{}); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func (tr TagRepository) AttachTagToArticle(attachBody *ArticleIdAndTagId) error {
	articleEntity := models.Article{ID: attachBody.ArticleId}
	tagEntity := models.Tag{ID: attachBody.TagId}

	return tr.Client.Model(&articleEntity).Association("Tags").Append(&tagEntity)
}

func (tr TagRepository) DetachTagFromArticle(detachBody *ArticleIdAndTagId) error {
	articleEntity := models.Article{ID: detachBody.ArticleId}
	tagEntity := models.Tag{ID: detachBody.TagId}

	return tr.Client.Model(&articleEntity).Association("Tags").Delete(&tagEntity)
}

func (tr TagRepository) DeleteTagAndItsAssociations(tagId uint) error {
	tagEntity := models.Tag{ID: tagId}

	if err := tr.Client.Model(&tagEntity).Association("Articles").Clear(); err != nil {
		return err
	}

	return tr.Client.Delete(&tagEntity, tagId).Error
}

func (tr TagRepository) GetArticlesByTagId(tagId uint) ([]*models.Article, error) {
	tagEntity := models.Tag{}

	err := tr.Client.Where("id = ?", tagId).Preload("Articles", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "created_at", "updated_at", "user_id", "title", "article_link", "og_image_link")
	}).First(&tagEntity).Error

	if err != nil {
		return nil, err
	}

	return tagEntity.Articles, nil
}

func (tr TagRepository) UpdateTagById(tagId uint, title string) error {
	tagEntity := models.Tag{ID: tagId}

	return tr.Client.Model(&tagEntity).Update("title", title).Error
}