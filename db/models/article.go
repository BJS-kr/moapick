package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Email       string `json:"email"`
	ArticleLink string `json:"article_link"`
	OgImageLink string `json:"ogimage_link"`
}
