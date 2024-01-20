package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Email    string    `gorm:"uniqueIndex:email_unique_idx" json:"email"`
	Articles []Article `json:"articles"`
	Tags	 []Tag	   `json:"tags"`
}
