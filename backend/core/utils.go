package core

import "gorm.io/gorm"

type Context struct {
	DB        *gorm.DB
	RootPath  string
	CachePath string
}
