package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Title    string    `json:"title"`
	User   	 User      `json:"user"`
	Articles []Article `json:"articles"`
}