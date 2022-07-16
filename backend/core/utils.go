package core

import "gorm.io/gorm"

type Context struct {
	DB        *gorm.DB
	RootPath  string
	CachePath string
}

// Files to ignore when processing media
var FilesToIgnore = []string{".DS_Store"}
