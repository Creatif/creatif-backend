package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"io"
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

func addToList(client *http.Client, projectId, name string, variable propertyVariable, references []map[string]string, imagePaths []string) httpResult {
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

	url := fmt.Sprintf("%s%s%s", URL, "/declarations/list/add/", projectId)
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
	response.Body.Close()

	if err != nil {
		return newHttpResult(nil, err, 0, false, Cannot_Continue_Procedure)
	}

	return newHttpResult(response, err, response.StatusCode, response.StatusCode >= 200 && response.StatusCode <= 299, Cannot_Continue_Procedure)
}

func addToMapAndGetAccountId(client *http.Client, projectId string, accountId string, account account) string {
	var genAccountId string
	handleHttpError(addToMap(client, projectId, accountId, account.variable, account.references, account.imagePaths), func(res *http.Response) error {
		if res.StatusCode < 200 || res.StatusCode > 299 {
			res.Body.Close()
			return errors.New(fmt.Sprintf("Generating one of the accounts return a status code %d", res.StatusCode))
		}

		b, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}

		genAccountId = m["id"].(string)

		return nil
	})

	return genAccountId
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
	images, err := generateBase64Images("images/profileImages", 1)
	if err != nil {
		printers["warning"].Printf("Unable to generate base64 image in Accounts generator %s\n", err.Error())
	}

	accountsToGenerate := 10

	accounts := make([]account, accountsToGenerate)

	for i := 0; i < accountsToGenerate; i++ {
		accountValueData := map[string]string{
			"name":       faker.FirstName(),
			"lastName":   faker.LastName(),
			"address":    faker.GetRealAddress().Address,
			"city":       faker.GetRealAddress().City,
			"postalCode": faker.GetRealAddress().PostalCode,
		}

		if len(images) == 1 {
			accountValueData["profileImage"] = images[0]
		}

		b, err := json.Marshal(accountValueData)
		if err != nil {
			return nil, err
		}

		uniqueName := uuid.New().String()
		acc := newAccount(uniqueName, nil, nil, newAccountVariable(uniqueName, "eng", "modifiable", "", string(b), pickRandomUniqueGroups(groupIds, 3)))
		accounts[i] = acc
	}

	return accounts, nil
}

func generatePropertiesStructureData(accountId string, groupIds []string) ([]property, error) {
	locales := []string{"eng", "afh", "kam", "ota", "oto"}

	properties := make([]property, 0)
	for _, locale := range locales {
		_, err := generateBase64Images("images/propertyImages", 3)
		if err != nil {
			printers["warning"].Printf("Unable to generate base64 image in Accounts generator %s\n", err.Error())
		}

		propertyStatutes := []string{"Rent", "Sell", "Rent business"}
		propertyTypes := []string{"House", "Apartment", "Studio", "Land"}

		for _, ps := range propertyStatutes {
			for _, pt := range propertyTypes {
				for i := 0; i < 10; i++ {
					p := generateSinglePropertyData(pt)
					//p["propertyImages"] = images
					p["propertyStatus"] = ps

					uniqueName := uuid.New().String()

					b, err := json.Marshal(p)
					if err != nil {
						return nil, err
					}

					properties = append(properties, newProperty(
						uniqueName,
						nil,
						nil,
						newPropertyVariable(
							uniqueName,
							locale,
							"modifiable",
							"",
							string(b),
							pickRandomUniqueGroups(groupIds, 3),
						),
					))
				}
			}
		}
	}

	return properties, nil
}

func generateSinglePropertyData(pt string) map[string]interface{} {
	propertyValueData := make(map[string]interface{})
	propertyValueData["address"] = faker.GetRealAddress().Address
	propertyValueData["city"] = faker.GetRealAddress().City
	propertyValueData["postalCode"] = faker.GetRealAddress().PostalCode

	if pt == "House" {
		propertyValueData["propertyType"] = "House"
		propertyValueData["numOfHouseFloors"] = randomBetween(1, 3)
		propertyValueData["houseSize"] = randomBetween(50, 500)
		propertyValueData["houseLocalPrice"] = randomBetween(1200, 5000)

		i := rand.IntN(10)
		if i%2 == 0 {
			propertyValueData["houseBackYard"] = true
			propertyValueData["houseBackYardSize"] = randomBetween(50, 500)
		} else {
			propertyValueData["houseBackYard"] = false
		}

		if i%5 == 0 {
			propertyValueData["houseNeedsRepair"] = true
			propertyValueData["houseRepairNote"] = faker.Sentence()
		} else {
			propertyValueData["houseNeedsRepair"] = false
		}
	}

	if pt == "Apartment" {
		propertyValueData["propertyType"] = "Apartment"
		propertyValueData["apartmentFloorNumber"] = randomBetween(10, 50)
		propertyValueData["apartmentSize"] = randomBetween(50, 500)
		propertyValueData["apartmentLocalPrice"] = randomBetween(500, 1500)

		i := rand.IntN(10)
		if i%2 == 0 {
			propertyValueData["apartmentBalcony"] = true
			propertyValueData["apartmentBalconySize"] = randomBetween(10, 30)
		} else {
			propertyValueData["apartmentBalcony"] = false
		}
	}

	if pt == "Studio" {
		propertyValueData["propertyType"] = "Studio"
		propertyValueData["studioFloorNumber"] = randomBetween(10, 50)
		propertyValueData["studioSize"] = randomBetween(20, 40)
	}

	if pt == "Land" {
		propertyValueData["propertyType"] = "Land"
		propertyValueData["landSize"] = randomBetween(1000, 5000)

		i := rand.IntN(10)
		if i%2 == 0 {
			propertyValueData["hasConstructionPermit"] = true
		} else {
			propertyValueData["hasConstructionPermit"] = false
		}
	}

	return propertyValueData
}
