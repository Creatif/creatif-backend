package main

import (
	"bufio"
	app2 "creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/logger"
	storage2 "creatif/pkg/lib/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
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
	if err := logger.BuildLoggers("/app/var/log"); err != nil {
		log.Fatalln(fmt.Sprintf("Cannot createProject logger: %s", err.Error()))
	}

	logger.Info("Health info logger health check... Ignore!")
	logger.Warn("Health warning logger health check... Ignore!")
	logger.Error("Health error logger health check... Ignore!")
}

func runPublic() {
	if _, err := os.Stat("/app/public"); os.IsNotExist(err) {
		err := os.MkdirAll("/app/public", os.ModePerm)

		if err != nil {
			log.Fatalln(fmt.Sprintf("Cannot create public directory: %s", err.Error()))
		}
	}
}

func runAssets() {
	if _, err := os.Stat("/app/var/assets"); os.IsNotExist(err) {
		err := os.MkdirAll("/app/var/assets", os.ModePerm)

		if err != nil {
			log.Fatalln(fmt.Sprintf("Cannot create public directory: %s", err.Error()))
		}
	}
}

func setupServer() *echo.Echo {
	srv := echo.New()

	srv.HideBanner = true
	srv.Server.ReadTimeout = 20 * time.Second
	srv.Server.WriteTimeout = 20 * time.Second

	return srv
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
	var exists declarations.Locale
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

	l := make([]declarations.Locale, 0)
	fileScanner.Scan()
	for fileScanner.Scan() {
		values := strings.Split(fileScanner.Text(), ",")
		l = append(l, declarations.NewLocale(values[3], values[0]))
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

func createDatabase() {
	ok, err := isMigrationPerformed()
	if err != nil {
		log.Fatalln(err)
	}

	if !ok {
		sqlDb := createSchemas()

		if _, err := sqlDb.Exec("ALTER DATABASE app SET search_path TO declarations;"); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.Group{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.File{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.VariableGroup{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.Map{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.MapVariable{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.List{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.ListVariable{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.Group{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.Reference{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(declarations.Locale{}); err != nil {
			log.Fatalln(err)
		}

		if _, err := sqlDb.Exec("ALTER DATABASE app SET search_path TO app;"); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(app2.Project{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(app2.User{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(app2.Event{}); err != nil {
			log.Fatalln(err)
		}

		if _, err := sqlDb.Exec("ALTER DATABASE app SET search_path TO published;"); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(published.PublishedList{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(published.PublishedFile{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(published.PublishedMap{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(published.PublishedReference{}); err != nil {
			log.Fatalln(err)
		}

		if err := storage2.Gorm().AutoMigrate(published.Version{}); err != nil {
			log.Fatalln(err)
		}
	}
}

func isMigrationPerformed() (bool, error) {
	// Most stupid check and migration ever
	var count int
	res := storage2.Gorm().Raw("SELECT count(schemaname) FROM pg_catalog.pg_tables where schemaname IN('declarations', 'app', 'published')").Scan(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count == 17, nil
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
	if _, err := sqlDb.Exec("CREATE SCHEMA IF NOT EXISTS published"); err != nil {
		log.Fatalln(err)
	}

	return sqlDb
}
