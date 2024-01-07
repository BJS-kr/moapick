package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	UserId 		uint   `gorm:"uniqueIndex:userid_title_unique_idx" json:"email"`
	User 		User   `gorm:"foreignKey:UserId" json:"user"`
	TagId		uint 	
	Tag 		Tag    `gorm:"foreignKey:TagId" json:"tag"`
	Title		string `gorm:"uniqueIndex:userid_title_unique_idx" json:"title"`			
	ArticleLink string `json:"article_link"`
	OgImageLink string `json:"og_image_link"`
}
