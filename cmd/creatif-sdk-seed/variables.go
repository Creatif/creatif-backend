package main

import (
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"net/http"
)

func addToMap(client *http.Client, projectId, name string, variable map[string]interface{}, references []map[string]string, imagePaths []string) httpResult {
	body := map[string]interface{}{
		"name":       name,
		"variable":   variable,
		"references": references,
		"imagePaths": imagePaths,
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

func addToList(client *http.Client, projectId, name string, variable map[string]interface{}, references []map[string]string, imagePaths []string) httpResult {
	body := map[string]interface{}{
		"name":       name,
		"variable":   variable,
		"references": references,
		"imagePaths": imagePaths,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/list/", projectId)
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

func preCreateAccountStructureData(structureName string) []map[string]interface{} {
	maps := make([]map[string]interface{}, 200)

	for {
		if len(maps) == 200 {
			return maps
		}

		m := make(map[string]interface{})
		name := fmt.Sprintf("%s-%s", faker.FirstName(), faker.LastName())

		isDuplicate := false
		for _, p := range maps {
			if p["name"] == name {
				isDuplicate = true
				break
			}
		}

		if isDuplicate {
			continue
		}

		m["name"] = structureName
		m["references"] = nil
		m["imagePaths"] = nil
		m["variable"] = map[string]interface{}{
			"name": name,
		}
	}
}
