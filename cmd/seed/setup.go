package main

import (
	"bufio"
	app2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/logger"
	storage2 "creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
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
		log.Fatalln(fmt.Sprintf("Cannot createProject logger: %s", err.Error()))
	}

	logger.Info("Health info logger health check... Ignore!")
	logger.Warn("Health warning logger health check... Ignore!")
	logger.Error("Health error logger health check... Ignore!")
}

func releaseAllLocks() error {
	var stat []int64
	res := storage2.Gorm().Raw(`SELECT DISTINCT pid FROM pg_locks l, pg_stat_all_tables t WHERE l.relation = t.relid AND t.relname = 'list_variables'`).Scan(&stat)
	if res.Error != nil {
		return res.Error
	}

	for s, _ := range stat {
		if res := storage2.Gorm().Exec("SELECT pg_cancel_backend(?)", s); res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func loadLocales() error {
	var exists app2.Locale
	if res := storage2.Gorm().First(&exists); res.Error != nil {
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return res.Error
		}
	}

	if exists.ID != "" {
		if err := locales.Store(); err != nil {
			return err
		}

		return nil
	}

	readFile, err := os.Open("/app/assets/locales.csv")
	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	l := make([]app2.Locale, 0)
	fileScanner.Scan()
	for fileScanner.Scan() {
		values := strings.Split(fileScanner.Text(), ",")
		l = append(l, app2.NewLocale(values[3], values[0]))
	}

	if err := readFile.Close(); err != nil {
		return err
	}

	if res := storage2.Gorm().Create(&l); res.Error != nil {
		return res.Error
	}

	if err := locales.Store(); err != nil {
		return err
	}

	return nil
}
