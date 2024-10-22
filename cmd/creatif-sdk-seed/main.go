package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
)

func main() {
	loadEnv()
	runDb()

	shouldJustCleanup := len(os.Args) == 2 && os.Args[1] == "cleanup"
	if shouldJustCleanup {
		doOrderedCleanup()
		os.Exit(0)
	}

	successColor := color.New(color.FgGreen).Add(color.Bold)

	anonymousClient := createAnonymousClient()
	email := "email@gmail.com"
	password := "password"

	if adminExists(anonymousClient).Ok() {
		printNewlineSandwich(successColor, "Admin already exists which means that the seed is there.\nIf it is not and this is a mistake, just delete the docker volume and try again.\nThis is fine and OK since this is a seed program to test the SDK.\nFeel free to abuse it.")
		return
	}

	printers["info"].Println("Creating admin and logging in...")
	handleError(createAdmin(anonymousClient, email, password), nil)

	authToken := extractAuthenticationCookie(handleError(login(anonymousClient, email, password), nil))

	authenticatedClient := createAuthenticatedClient(authToken)

	printers["info"].Println("Creating projects...")

	projectNames := []string{"Warsaw Brokers", "London Brokers", "Paris Brokers", "Berlin Brokers", "Barcelona Brokers"}
	projects := make([]map[string]string, len(projectNames))
	for i, p := range projectNames {
		handleError(createProject(authenticatedClient, p), func(res *http.Response) error {
			var m map[string]interface{}
			b, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(b, &m); err != nil {
				return err
			}

			// interface conversion is ok here since I know that it will be a string
			projects[i] = map[string]string{
				"id":   m["id"].(string),
				"name": m["name"].(string),
			}

			return nil
		})
	}

	fmt.Println(projects)
}
