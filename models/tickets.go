package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`

	UserID uint `json:"user_id"`
}