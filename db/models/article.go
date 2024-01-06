package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Email       string `gorm:"index:email_title_unique_idx,unique" json:"email"`
	Title			 	string `gorm:"index:email_title_unique_idx,unique" json:"title"`
	ArticleLink string `json:"article_link"`
	OgImageLink string `json:"ogimage_link"`
}
