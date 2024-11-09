package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var gormHandle *gorm.DB

func Connect(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return err
	}

	d, err := db.DB()
	if err != nil {
		return err
	}

	if err := d.Ping(); err != nil {
		return err
	}

	d.SetMaxIdleConns(100)
	d.SetMaxOpenConns(10000)
	d.SetConnMaxLifetime(time.Hour)

	gormHandle = db

	return nil
}
