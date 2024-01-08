package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	User 		User   `json:"user"`
	Tag 		Tag    `json:"tag"`
	Title		string `gorm:"uniqueIndex:userid_title_unique_idx" json:"title"`			
	ArticleLink string `json:"article_link"`
	OgImageLink string `json:"og_image_link"`
}
