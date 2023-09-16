package main

import (
	"creatif/pkg/app/declarations/createVariable"
	storage2 "creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	loadEnv()
	db()

	createDeclarationVariablesWithoutValue(1000)
}

func createDeclarationVariablesWithoutValue(num int) {
	for i := 0; i < num; i++ {
		b, _ := json.Marshal(map[string]interface{}{
			"one":   "one",
			"two":   []string{"group1", "group2", "group3"},
			"three": []int{1, 2, 3, 4},
			"four":  4,
		})

		handler := createVariable.New(createVariable.NewModel(
			fmt.Sprintf("name-%d", i),
			"modifiable",
			[]string{"one", "two", "three"},
			b,
			b,
		))

		_, err := handler.Handle()
		if err != nil {
			log.Fatal(err)
		}
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
