package main

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	storage2 "creatif/pkg/lib/storage"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	loadEnv()
	db()
	runMigrations()
	closeConnection()
}

func closeConnection() {
	s, err := storage2.SQLDB()
	if err != nil {
		log.Fatalln(err)
	}

	if err := s.Close(); err != nil {
		log.Fatalln(err)
	}
}

func runMigrations() {
	sqlDb := createSchemas()

	if _, err := sqlDb.Exec("ALTER DATABASE app SET search_path TO declarations;"); err != nil {
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.Variable{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.Map{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.MapVariable{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.List{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.ListVariable{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.Group{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.Reference{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.Locale{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if _, err := sqlDb.Exec("ALTER DATABASE app SET search_path TO app;"); err != nil {
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(app.Project{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(app.User{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}
	/*
		storage2.Gorm().Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT unique_variable_per_project UNIQUE (name, project_id)", domain.VARIABLES_TABLE))
		storage2.Gorm().Exec(fmt.Sprintf("CREATE UNIQUE INDEX idx_unique_variable_name ON %s (name) WHERE project_id IS NOT NULL", domain.VARIABLES_TABLE))*/
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func db() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Zagreb",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	err := storage2.Connect(dsn)

	if err != nil {
		log.Fatalln(err)
	}
}

func createSchemas() *sql.DB {
	sqlDb, err := storage2.SQLDB()
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := sqlDb.Exec("CREATE SCHEMA IF NOT EXISTS app"); err != nil {
		log.Fatalln(err)
	}

	if _, err := sqlDb.Exec("CREATE SCHEMA IF NOT EXISTS declarations"); err != nil {
		log.Fatalln(err)
	}

	if _, err := sqlDb.Exec("CREATE SCHEMA IF NOT EXISTS assignments"); err != nil {
		log.Fatalln(err)
	}

	if _, err := sqlDb.Exec("CREATE SCHEMA IF NOT EXISTS content"); err != nil {
		log.Fatalln(err)
	}

	return sqlDb
}
