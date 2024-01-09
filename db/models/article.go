package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	UserId  	  uint   `gorm:"uniqueIndex:userid_title_unique_idx" json:"user_id"`
	Tags 		    []*Tag `gorm:"many2many:article_tags;" json:"tags"`
	Title		    string `gorm:"uniqueIndex:userid_title_unique_idx" json:"title"`			
	ArticleLink string `json:"article_link"`
	OgImageLink string `json:"og_image_link"`
}
