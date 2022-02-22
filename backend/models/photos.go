package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Path string `gorm:"not null"`
}
