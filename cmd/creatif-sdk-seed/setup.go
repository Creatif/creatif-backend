package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
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

	err := Connect(dsn)

	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot connect to database: %s", err.Error()))
	}
}

func reactToFlags() {
	shouldJustCleanup := len(os.Args) == 2 && os.Args[1] == "--cleanup"
	if shouldJustCleanup {
		doOrderedCleanup()
		os.Exit(0)
	}
}

func preSeedAuthAndSetup(client *http.Client) *http.Client {
	if adminExists(client).Ok() {
		printNewlineSandwich(printers["success"], "Admin already exists which means that the seed is there.\nIf it is not and this is a mistake, just delete the docker volume and try again.\nThis is fine and OK since this is a seed program to test the SDK.\nFeel free to abuse it.")
		return nil
	}

	printers["info"].Println("Creating admin and logging in")
	handleHttpError(createAdmin(client, Email, Password), nil)

	authToken := extractAuthenticationCookie(handleHttpError(login(client, Email, Password), nil))

	return createAuthenticatedClient(authToken)
}
