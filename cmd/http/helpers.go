package main

import (
	"creatif/pkg/lib/logger"
	storage2 "creatif/pkg/lib/storage"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"time"
)

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func runDb() {
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
		log.Fatalln(fmt.Sprintf("Cannot connect to database: %s", err.Error()))
	}
}

func runLogger() {
	if err := logger.BuildLoggers(os.Getenv("LOG_DIRECTORY")); err != nil {
		log.Fatalln(fmt.Sprintf("Cannot createNode logger: %s", err.Error()))
	}

	logger.Info("Health info logger health check... Ignore!")
	logger.Warn("Health warning logger health check... Ignore!")
	logger.Error("Health error logger health check... Ignore!")
}

func runAssets() {
	assetsDir := os.Getenv("ASSETS_DIRECTORY")
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		err := os.MkdirAll(assetsDir, os.ModePerm)

		if err != nil {
			log.Fatalln(fmt.Sprintf("Cannot createNode assets directory: %s", err.Error()))
		}
	}
}

func setupServer() *echo.Echo {
	srv := echo.New()

	srv.Server.ReadTimeout = 20 * time.Second
	srv.Server.WriteTimeout = 20 * time.Second

	return srv
}
