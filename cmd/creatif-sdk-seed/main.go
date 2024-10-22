package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	loadEnv()
	runDb()

	client := newClient(newClientParams(&http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost:     1024,
		TLSHandshakeTimeout: 0,
	}, nil, nil, 0))

	email := "email@gmail.com"
	password := "password"

	handleError(createAdmin(client, email, password))

	loginResult := login(client, email, password)
	fmt.Println(loginResult)

	/*	fmt.Println("created admin")
		projects := []string{"Warsaw Brokers", "London Brokers", "Paris Brokers", "Berlin Brokers", "Barcelona Brokers"}
		projectIds := make([]string, len(projects))
		for i, p := range projects {
			projectIds[i] = createProject(p, login())
		}
		fmt.Println("created projects")

		for _, p := range projectIds {
			createAccountStructure(p, "Accounts")
			createPropertiesStructure(p, "Properties")
		}
	*/
}
