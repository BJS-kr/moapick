package tag

import (
	"errors"
	"moapick/db"
	"moapick/db/models"

	"gorm.io/gorm"
)

func CreateTag(title string, userId uint) error {
	tagEntity := models.Tag{
		Title:  title,
		UserID: userId,
	}

	result := db.Client.Create(&tagEntity)

	return result.Error
}

func GetAllTagsOfUser(userId uint) ([]models.Tag, error) {
	var tags []models.Tag
	result := db.Client.Where("user_id = ?", userId).Find(&tags)

	return tags, result.Error
}

func IsTagBelongsToUser(userId, tagId uint) (bool, error) {
	if result := db.Client.Where("ID = ? AND user_id = ?", tagId, userId).First(&models.Tag{}); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func AttachTagToArticle(attachBody *AttachBody) error {
	articleEntity := models.Article{
		Model: gorm.Model{
			ID: attachBody.ArticleId,
		},
	}

	return db.Client.Model(&articleEntity).Association("Tags").Append(&models.Tag{ID: attachBody.TagId})
}
