package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Title    string    `json:"title"`
	UserID   uint      `json:"user_id"`
	Articles []*Article `gorm:"many2many:article_tags;" json:"articles"`
}