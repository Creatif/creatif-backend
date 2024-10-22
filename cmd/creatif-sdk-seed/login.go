package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func login(client *http.Client, email, password string) httpResult {
	body := map[string]string{
		"email":    email,
		"password": password,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s", URL, "/app/auth/login")
	req, err := newRequest(request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Url:    url,
		Method: "POST",
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

func extractAuthenticationCookie(httpResult httpResult) string {
	response := httpResult.Response()
	cookies := response.Cookies()

	for _, c := range cookies {
		if c.Name == "api_authentication" {
			return c.Value
		}
	}

	return ""
}
