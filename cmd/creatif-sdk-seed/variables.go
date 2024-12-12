package main

import (
	"creatif-sdk-seed/dataGeneration"
	"creatif-sdk-seed/errorHandler"
	http2 "creatif-sdk-seed/http"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func createClient(client *http.Client, projectId, name string, variable dataGeneration.ClientVariable, connections []map[string]string, imagePaths []string) http2.HttpResult {
	body := map[string]interface{}{
		"name": name,
		"variable": map[string]interface{}{
			"name":      variable.Name,
			"locale":    variable.Locale,
			"behaviour": variable.Behaviour,
			"groups":    variable.Groups,
			"metadata":  "",
			"value":     variable.Value,
		},
		"connections": connections,
		"imagePaths":  imagePaths,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/map/add/", projectId)
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

func createManager(client *http.Client, mapId, projectId string, variable dataGeneration.ManagerVariable, connections []map[string]string, imagePaths []string) http2.HttpResult {
	body := map[string]interface{}{
		"name": mapId,
		"variable": map[string]interface{}{
			"name":      variable.Name,
			"locale":    variable.Locale,
			"behaviour": variable.Behaviour,
			"groups":    variable.Groups,
			"metadata":  "",
			"value":     variable.Value,
		},
		"connections": connections,
		"imagePaths":  imagePaths,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/map/add/", projectId)
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

func createProperty(client *http.Client, projectId, name string, variable dataGeneration.PropertyVariable, connections []map[string]string, imagePaths []string) http2.HttpResult {
	body := map[string]interface{}{
		"name": name,
		"variable": map[string]interface{}{
			"name":      variable.Name,
			"locale":    variable.Locale,
			"behaviour": variable.Behaviour,
			"groups":    variable.Groups,
			"metadata":  "",
			"value":     variable.Value,
		},
		"connections": connections,
		"imagePaths":  imagePaths,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/list/add/", projectId)
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
	if response != nil && response.Body != nil {
		response.Body.Close()
	}

	if err != nil {
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	return http2.NewHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299, Cannot_Continue_Procedure)
}

func addToMapAndGetClientId(client *http.Client, projectId string, clientId string, clientVariable dataGeneration.Client) string {
	var getClientId string
	httpResult := errorHandler.HandleHttpError(createClient(client, projectId, clientId, clientVariable.Variable, clientVariable.Connections, clientVariable.ImagePaths))
	res := httpResult.Response()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		if res.Body != nil {
			res.Body.Close()
		}
		errorHandler.HandleAppError(errors.New(fmt.Sprintf("Generating one of the clients return a status code %d", res.StatusCode)), Cannot_Continue_Procedure)
	}

	if res.Body == nil {
		errorHandler.HandleAppError(errors.New("addToMapAndGetClientId() was trying to get a response body on a nil body"), Cannot_Continue_Procedure)
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
	}

	getClientId = m["id"].(string)

	return getClientId
}
