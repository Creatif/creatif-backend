package main

import (
	"database/sql"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	return gormHandle
}

func SQLDB() (*sql.DB, error) {
	s, err := Gorm().DB()
	if err != nil {
		return nil, err
	}

	return s, nil
}
