package main

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"io"
	"net/http"
)

func createAdmin(client *http.Client, email, password string) httpResult {
	body := map[string]string{
		"name":     faker.Name(),
		"lastName": faker.LastName(),
		"email":    email,
		"password": password,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s", URL, "/app/auth/admin/create")
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

	return newHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299, Cannot_Continue_Procedure)
}

func adminExists(client *http.Client) httpResult {
	url := fmt.Sprintf("%s%s", URL, "/app/auth/admin/exists")
	req, err := newRequest(request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Url:    url,
		Method: "GET",
		Body:   nil,
	})
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	response, err := Make(req, client)

	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	b, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}
	var doesExist bool
	err = json.Unmarshal(b, &doesExist)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	return newHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299 && doesExist, Can_Continue)
}
