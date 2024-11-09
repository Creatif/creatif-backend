package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var gormHandle *gorm.DB

type postgresDb struct {
	db *gorm.DB
}

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

	d.SetMaxIdleConns(5)
	d.SetMaxOpenConns(10)
	d.SetConnMaxLifetime(30 * time.Minute) // Set the maximum connection lifetime

	if err := initShortId(); err != nil {
		return err
	}

	gormHandle = db

	return nil
}
