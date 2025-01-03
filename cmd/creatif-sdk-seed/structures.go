package main

import (
	"creatif-sdk-seed/errorHandler"
	http2 "creatif-sdk-seed/http"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func createMapStructure(client *http.Client, projectId, name string) http2.HttpResult {
	body := map[string]string{
		"name": name,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/map/", projectId)
	req, err := http2.NewRequest(http2.Request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Url:    url,
		Method: "PUT",
		Body:   b,
	})
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	response, err := http2.Make(req, client)

	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	return http2.NewHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299, Cannot_Continue_Procedure)
}

func createListStructure(client *http.Client, projectId, name string) http2.HttpResult {
	body := map[string]string{
		"name": name,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/list/", projectId)
	req, err := http2.NewRequest(http2.Request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Url:    url,
		Method: "PUT",
		Body:   b,
	})
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	response, err := http2.Make(req, client)

	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	return http2.NewHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299, Cannot_Continue_Procedure)
}

func createClientStructureAndReturnId(client *http.Client, projectId string) string {
	var id string
	result := errorHandler.HandleHttpError(createMapStructure(client, projectId, "Clients"))
	res := result.Response()

	if res.Body == nil {
		errorHandler.HandleAppError(errors.New("createPropertiesStructureAndReturnID() does not have a body"), Cannot_Continue_Procedure)
	}

	defer res.Body.Close()
	var m map[string]interface{}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	if err := json.Unmarshal(b, &m); err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	id = m["id"].(string)

	return id
}

func createManagersStructureAndReturnId(client *http.Client, projectId string) string {
	var id string
	result := errorHandler.HandleHttpError(createMapStructure(client, projectId, "Managers"))
	res := result.Response()

	if res.Body == nil {
		errorHandler.HandleAppError(errors.New("createManagersStructureAndReturnId() does not have a body"), Cannot_Continue_Procedure)
	}

	defer res.Body.Close()
	var m map[string]interface{}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	if err := json.Unmarshal(b, &m); err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	id = m["id"].(string)

	return id
}

func createPropertiesStructureAndReturnID(client *http.Client, projectId string) string {
	var id string
	result := errorHandler.HandleHttpError(createListStructure(client, projectId, "Properties"))
	res := result.Response()

	if res.Body == nil {
		errorHandler.HandleAppError(errors.New("createPropertiesStructureAndReturnID() does not have a body"), Cannot_Continue_Procedure)
	}

	defer res.Body.Close()
	var m map[string]interface{}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	if err := json.Unmarshal(b, &m); err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	id = m["id"].(string)

	return id
}
