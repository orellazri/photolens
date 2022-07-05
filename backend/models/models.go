package models

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	Path         string `gorm:"not null"`
	LastModified time.Time
}
