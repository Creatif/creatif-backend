package main

import (
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	return gormHandle
}
