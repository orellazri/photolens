package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Media struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Path         string `gorm:"not null"`
	ContentType  string `gorm:"not null"`
	IsPhoto      bool
	LastModified time.Time
}

type Meta struct {
	Key   string `gorm:"primaryKey;not null"`
	Value string `gorm:"not null"`
}
