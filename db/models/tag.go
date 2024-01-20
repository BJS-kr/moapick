package models

import "time"

// tag를 soft delete할 필요는 없다고 판단
type Tag struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Title     string     `gorm:"uniqueIndex:userid_title_uniq_idx" json:"title"`
	UserID    uint       `gorm:"uniqueIndex:userid_title_uniq_idx" json:"user_id,omitempty"`
	Articles  []*Article `gorm:"many2many:article_tags;" json:"articles,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
