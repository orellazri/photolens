package models

import (
	"time"
)

type Media struct {
	Path         string `gorm:"primaryKey"`
	ContentType  string `gorm:"not null"`
	IsPhoto      bool
	LastModified time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
