package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Title    string     `gorm:"uniqueIndex:tag_title_uniq_idx" json:"title"`
	UserID   uint       `json:"user_id"`
	Articles []*Article `gorm:"many2many:article_tags;" json:"articles"`
}
