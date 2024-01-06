package article

import (
	"fmt"
	"moapick/db/models"
	"net/url"
)

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func SaveArticle(userEmail string, article ArticleBody) (*models.Article, error) {
	if !isValidURL(article.Link) {
		return nil, fmt.Errorf("requested article link is not valid url")
	}

	savedArticle, err := SaveArticleEntity(userEmail, article.Link)

	if err != nil {
		return nil, err
	}

	return savedArticle, nil
}
