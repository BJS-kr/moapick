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
	if result := db.Client.Select("id").Where("id = ? AND user_id = ?", tagId, userId).First(&models.Tag{}); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

func AttachTagToArticle(attachBody *ArticleIdAndTagId) error {
	articleEntity := models.Article{ID: attachBody.ArticleId}
	tagEntity := models.Tag{ID: attachBody.TagId}

	return db.Client.Model(&articleEntity).Association("Tags").Append(&tagEntity)
}

func DetachTagFromArticle(detachBody *ArticleIdAndTagId) error {
	articleEntity := models.Article{ID: detachBody.ArticleId}
	tagEntity := models.Tag{ID: detachBody.TagId}

	return db.Client.Model(&articleEntity).Association("Tags").Delete(&tagEntity)
}

func DeleteTagAndItsAssociations(tagId uint) error {
	tagEntity := models.Tag{ID: tagId}

	if err := db.Client.Model(&tagEntity).Association("Articles").Clear(); err != nil {
		return err
	}

	return db.Client.Delete(&tagEntity, tagId).Error
}
