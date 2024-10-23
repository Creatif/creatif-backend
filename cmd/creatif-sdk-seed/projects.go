package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func createProject(client *http.Client, name string) httpResult {
	body := map[string]string{
		"name": name,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s", URL, "/app/project")
	req, err := newRequest(request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Url:    url,
		Method: "PUT",
		Body:   b,
	})
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	response, err := Make(req, client)

	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	return newHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299, Can_Continue)
}

func generateProjects(client *http.Client) []project {
	projectNames := []string{"Warsaw Brokers", "London Brokers", "Paris Brokers", "Berlin Brokers", "Barcelona Brokers"}
	projects := make([]project, len(projectNames))
	for i, p := range projectNames {
		handleHttpError(createProject(client, p), func(res *http.Response) error {
			var m project
			b, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(b, &m); err != nil {
				return err
			}

			projects[i] = m

			return nil
		})
	}

	return projects
}
