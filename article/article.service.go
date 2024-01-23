package article

import (
	"net/url"
)

type ArticleService struct{}

func (as ArticleService) IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}
