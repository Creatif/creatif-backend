package dataGeneration

import (
	"encoding/json"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"math/rand/v2"
)

type PropertyVariable struct {
	Name      string
	Locale    string
	Behaviour string
	Groups    []string
	Metadata  string
	Value     string
}

type Property struct {
	Name        string
	Connections []map[string]string
	ImagePaths  []string
	Variable    PropertyVariable
}

func GenerateSingleProperty(clientId, locale, propertyStatus, propertyType string, groupIds []string) (Property, error) {
	images, err := generateBase64Images("images/propertyImages", 3)
	if err != nil {
		printers["warning"].Printf("Unable to generate base64 image in Clients generator %s\n", err.Error())
	}

	p := generateSinglePropertyVariable(propertyType)
	p["propertyImages"] = images
	p["propertyStatus"] = propertyStatus

	uniqueName := uuid.New().String()

	b, err := json.Marshal(p)
	if err != nil {
		return Property{}, err
	}

	return newProperty(
		uniqueName,
		[]map[string]string{
			{
				"name":          "clients",
				"structureType": "map",
				"variableId":    clientId,
			},
		},
		[]string{"propertyImages"},
		newPropertyVariable(
			uniqueName,
			locale,
			"modifiable",
			"",
			string(b),
			pickRandomUniqueGroups(groupIds, 3),
		),
	), nil
}

func generateSinglePropertyVariable(pt string) map[string]interface{} {
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

func newProperty(name string, connections []map[string]string, imagePaths []string, variable PropertyVariable) Property {
	return Property{
		Name:        name,
		Connections: connections,
		ImagePaths:  imagePaths,
		Variable:    variable,
	}
}

func newPropertyVariable(name, locale, behaviour, metadata, value string, groups []string) PropertyVariable {
	return PropertyVariable{
		Name:      name,
		Locale:    locale,
		Behaviour: behaviour,
		Groups:    groups,
		Metadata:  metadata,
		Value:     value,
	}
}
