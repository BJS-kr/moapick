package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	UserId uint   `json:"user_id"`
	User   User   `gorm:"foreignKey:UserId" json:"user"`
	Title  string `json:"title"`
}