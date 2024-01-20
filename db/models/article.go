package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserId  	uint   `gorm:"uniqueIndex:userid_title_unique_idx" json:"user_id,omitempty"`
	Tags 		[]*Tag `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	Title		string `gorm:"uniqueIndex:userid_title_unique_idx" json:"title"`			
	ArticleLink string `json:"article_link"`
	OgImageLink string `json:"og_image_link"`
}
