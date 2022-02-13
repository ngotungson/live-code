package models

import "time"

type CommentPayload struct {
	Content string `json:"content"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `json:"-"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
