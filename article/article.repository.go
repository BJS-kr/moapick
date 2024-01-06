package article

import (
	"moapick/db"
	"moapick/db/models"

	"github.com/otiai10/opengraph"
)

func SaveArticleEntity(email, articleLink string) (*models.Article, error) {

	articleEntity := models.Article{
		Email:       email,
		ArticleLink: articleLink,
	}

	og, err := opengraph.Fetch(articleLink)

	if err == nil {
		if len(og.Image) > 0 {
			articleEntity.OgImageLink = og.Image[0].URL
		}
	}

	// TODO email + article link가 unique하게 만들기
	result := db.Client.Create(&articleEntity)

	if result.Error != nil {
		return nil, result.Error
	}

	return &articleEntity, nil
}
