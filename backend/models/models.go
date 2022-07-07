package models

import (
	"time"

	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	Path         string `gorm:"not null"`
	ContentType  string `gorm:"not null"`
	IsPhoto      bool
	LastModified time.Time
}
