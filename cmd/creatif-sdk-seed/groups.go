package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func createGroups(client *http.Client, projectId string) httpResult {
	groups := make([]map[string]string, 100)
	for i := 0; i < 100; i++ {
		groups[i] = map[string]string{
			"id":     "",
			"name":   fmt.Sprintf("group-%d", i),
			"action": "create",
			"type":   "",
		}
	}

	body := map[string]interface{}{
		"groups": groups,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s", URL, "/app/groups/"+projectId)
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

func createGroupsAndGetGroupIds(client *http.Client, projectId string) []string {
	groupIds := make([]string, 0)
	handleHttpError(createGroups(client, projectId), func(res *http.Response) error {
		b, _ := io.ReadAll(res.Body)
		var groups []map[string]string
		if err := json.Unmarshal(b, &groups); err != nil {
			return err
		}

		if err := res.Body.Close(); err != nil {
			return err
		}

		for _, g := range groups {
			groupIds = append(groupIds, g["id"])
		}

		return nil
	})

	return groupIds
}
