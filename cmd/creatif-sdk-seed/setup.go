package main

import (
	"creatif-sdk-seed/errorHandler"
	"creatif-sdk-seed/storage"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

	err := storage.Connect(dsn)

	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot connect to database: %s", err.Error()))
	}
}

func doOperations(operations []string) {
	for _, op := range operations {
		if op == "--regenerate" {
			doOrderedCleanup()
		}
	}
}

func processFlags() (int, int, error) {
	defaultNumberOfProjects := 1
	defaultPropertiesPerType := 1
	if len(os.Args) == 1 {
		return defaultNumberOfProjects, defaultNumberOfProjects, nil
	}

	// validation first
	for i := 0; i < len(os.Args[0:]); i++ {
		re := regexp.MustCompile(`^--projects=(\d+)$`)
		matches := re.FindStringSubmatch(os.Args[i])
		if len(matches) > 1 {
			numOfProjects, err := strconv.ParseInt(matches[1], 10, 32)
			if err != nil {
				return 0, 0, err
			}

			if numOfProjects < 1 || numOfProjects > 10 {
				return 0, 0, errors.New("Number of projects must be minimal 1 and below 10")
			}

			defaultNumberOfProjects = int(numOfProjects)
		}

		re = regexp.MustCompile(`^--properties-per-type=(\d+)$`)
		matches = re.FindStringSubmatch(os.Args[i])
		if len(matches) > 1 {
			propertiesPerStatus, err := strconv.ParseInt(matches[1], 10, 32)
			if err != nil {
				return 0, 0, err
			}

			if propertiesPerStatus < 1 || propertiesPerStatus > 10 {
				return 0, 0, errors.New("Number of properties per type must be minimal 1 and below 10")
			}

			defaultPropertiesPerType = int(propertiesPerStatus)
		}
	}

	operations := []string{}
	for i := 0; i < len(os.Args[0:]); i++ {
		// just cleanup the system and exit
		shouldJustCleanup := os.Args[i] == "--cleanup"
		if shouldJustCleanup {
			doOrderedCleanup()
			os.Exit(0)
		}

		shouldJustOutputHelp := os.Args[i] == "--help"
		if shouldJustOutputHelp {
			fmt.Print(help)
			os.Exit(0)
		}

		// just cleanup the system and start all over without existing
		shouldRegenerate := os.Args[i] == "--regenerate"
		if shouldRegenerate {
			operations = append(operations, "--regenerate")
		}

		optionalArgs := strings.Split(os.Args[i], "=")
		// if there are no optional flags, just continue
		if len(optionalArgs) == 0 {
			return defaultNumberOfProjects, defaultPropertiesPerType, nil
		}
	}

	doOperations(operations)
	return defaultNumberOfProjects, defaultPropertiesPerType, nil
}

func preSeedAuthAndSetup(client *http.Client) *http.Client {
	if adminExists(client).Ok() {
		printNewlineSandwich(printers["success"], "Admin already exists which means that the seed is there.\nIf it is not and this is a mistake, just delete the docker volume and try again.\nThis is fine and OK since this is a seed program to test the SDK.\nFeel free to abuse it.")
		return nil
	}

	printers["info"].Println("Creating admin and logging in")
	errorHandler.HandleHttpError(createAdmin(client, Email, Password))

	authToken := extractAuthenticationCookie(errorHandler.HandleHttpError(login(client, Email, Password)))

	return createAuthenticatedClient(authToken)
}
