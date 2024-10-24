package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func createMapStructure(client *http.Client, projectId, name string) httpResult {
	body := map[string]string{
		"name": name,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/map/", projectId)
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

func createListStructure(client *http.Client, projectId, name string) httpResult {
	body := map[string]string{
		"name": name,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/map/", projectId)
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

func createAccountStructureAndReturnID(client *http.Client, projectId string) string {
	var id string
	handleHttpError(createMapStructure(client, projectId, "Accounts"), func(res *http.Response) error {
		defer res.Body.Close()
		var m map[string]interface{}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}

		id = m["id"].(string)

		return nil
	})

	return id
}

func createPropertiesStructureAndReturnID(client *http.Client, projectId string) string {
	var id string
	handleHttpError(createListStructure(client, projectId, "Properties"), func(res *http.Response) error {
		defer res.Body.Close()
		var m map[string]interface{}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}

		id = m["id"].(string)

		return nil
	})

	return id
}
