package main

import (
	"creatif/pkg/app/domain/assignments"
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
	sql, err := storage2.SQLDB()
	if err != nil {
		log.Fatalln(err)
	}

	if err := sql.Close(); err != nil {
		log.Fatalln(err)
	}
}

func runMigrations() {
	sqlDb := createSchemas()

	if _, err := sqlDb.Exec("ALTER DATABASE app SET search_path TO declarations,app,assignments,content;"); err != nil {
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.Node{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(declarations.Map{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if _, err := sqlDb.Exec("ALTER DATABASE app SET search_path TO assignments,app,declarations,content;"); err != nil {
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(assignments.Node{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(assignments.NodeText{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(assignments.NodeBoolean{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}
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

	if _, err := sqlDb.Exec("DROP SCHEMA IF EXISTS public"); err != nil {
		log.Fatalln(err)
	}

	return sqlDb
}
