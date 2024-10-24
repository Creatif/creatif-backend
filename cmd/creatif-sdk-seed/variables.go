package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"math/rand/v2"
	"net/http"
)

func addToMap(client *http.Client, projectId, name string, variable accountVariable, references []map[string]string, imagePaths []string) httpResult {
	body := map[string]interface{}{
		"name": name,
		"variable": map[string]interface{}{
			"name":      variable.name,
			"locale":    variable.locale,
			"behaviour": variable.behaviour,
			"groups":    variable.groups,
			"metadata":  "",
			"value":     variable.value,
		},
		"references": references,
		"imagePaths": imagePaths,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/map/add/", projectId)
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

func addToList(client *http.Client, projectId, name string, variable map[string]interface{}, references []map[string]interface{}, imagePaths []string) httpResult {
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

func pickRandomUniqueGroups(groupIds []string, numOfGroups int) []string {
	picked := 0
	groups := make([]string, 3)
	for {
		if picked == numOfGroups {
			return groups
		}

		g := groupIds[rand.IntN(100)]

		isDuplicate := false
		for _, pickedGroup := range groups {
			if pickedGroup == g {
				isDuplicate = true
			}
		}

		if !isDuplicate {
			groups[picked] = g
			picked++
		}
	}
}

func generateAccountStructureData(groupIds []string) ([]account, error) {
	_, err := generateBase64Images("images/profileImages", 1)
	if err != nil {
		printers["warning"].Printf("Unable to generate base64 image in Accounts generator %s\n", err.Error())
	}

	accounts := make([]account, 200)
	successIterations := 0
	for {
		if successIterations == 200 {
			return accounts, nil
		}

		firstName := faker.FirstName()
		lastName := faker.LastName()
		name := fmt.Sprintf("%s-%s", firstName, lastName)

		isDuplicate := false
		for _, p := range accounts {
			if p.name == name {
				isDuplicate = true
				break
			}
		}

		if isDuplicate {
			continue
		}

		accountValueData := map[string]string{
			"name":       firstName,
			"lastName":   lastName,
			"address":    faker.GetRealAddress().Address,
			"city":       faker.GetRealAddress().City,
			"postalCode": faker.GetRealAddress().PostalCode,
		}
		/*
			if len(images) == 1 {
				accountValueData["profileImage"] = images[0]
			}*/

		b, err := json.Marshal(accountValueData)
		if err != nil {
			return nil, err
		}

		acc := newAccount(name, nil, nil, newAccountVariable(name, "eng", "modifiable", "", string(b), pickRandomUniqueGroups(groupIds, 3)))

		accounts[successIterations] = acc
		successIterations += 1
	}
}
