package main

import (
	"creatif/pkg/app/domain"
	storage2 "creatif/pkg/lib/storage"
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
	if err := storage2.Gorm().AutoMigrate(domain.Project{}); err != nil {
		closeConnection()
		log.Fatalln(err)
	}

	if err := storage2.Gorm().AutoMigrate(domain.User{}); err != nil {
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
