package models

import "time"

// tag를 soft delete할 필요는 없다고 판단
type Tag struct {
	ID        uint       `gorm:"primarykey"`
	Title     string     `gorm:"uniqueIndex:userid_title_uniq_idx" json:"title"`
	UserID    uint       `gorm:"uniqueIndex:userid_title_uniq_idx" json:"user_id"`
	Articles  []*Article `gorm:"many2many:article_tags;" json:"articles"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
