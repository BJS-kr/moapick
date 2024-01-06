package article

import (
	"moapick/db"
	"moapick/db/models"
)

func SaveArticle(articleEntity *models.Article) error {
	result := db.Client.Create(articleEntity)

	return result.Error
}
