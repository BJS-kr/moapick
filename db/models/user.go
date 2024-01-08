package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string    `gorm:"uniqueIndex:email_unique_idx" json:"email"`
	Articles []Article `json:"articles"`
	Tags	 []Tag	   `json:"tags"`
}
