package main

import (
	http2 "creatif-sdk-seed/http"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func createGroups(client *http.Client, projectId string) http2.HttpResult {
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
		return http2.NewHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s", URL, "/app/groups/"+projectId)
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

func createGroupsAndGetGroupIds(client *http.Client, projectId string) []string {
	groupIds := make([]string, 0)
	result := handleHttpError(createGroups(client, projectId))
	res := result.Response()

	if res == nil || res.Body == nil {
		handleAppError(errors.New("createGroupsAndGetGroupIds() is trying to work on nil body"), Cannot_Continue_Procedure)
	}

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	var groups []map[string]string
	if err := json.Unmarshal(b, &groups); err != nil {
		handleAppError(err, Cannot_Continue_Procedure)
	}

	if err := res.Body.Close(); err != nil {
		handleAppError(err, Cannot_Continue_Procedure)
	}

	for _, g := range groups {
		groupIds = append(groupIds, g["id"])
	}

	return groupIds
}
