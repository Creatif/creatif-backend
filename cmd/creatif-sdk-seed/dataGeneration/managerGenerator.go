package dataGeneration

import (
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type ManagerVariable struct {
	Name      string
	Locale    string
	Behaviour string
	Groups    []string
	Metadata  string
	Value     string
}

type Manager struct {
	Name        string
	Connections []map[string]string
	ImagePaths  []string
	Variable    ManagerVariable
}

type ClientConnection struct {
	StructureType string
	VariableID    string
}

func newManagerVariable(name, locale, behaviour, metadata, value string, groups []string) ManagerVariable {
	return ManagerVariable{
		Name:      name,
		Locale:    locale,
		Behaviour: behaviour,
		Groups:    groups,
		Metadata:  metadata,
		Value:     value,
	}
}

func newManager(name string, connections []map[string]string, imagePaths []string, variable ManagerVariable) Manager {
	return Manager{
		Name:        name,
		Connections: connections,
		ImagePaths:  imagePaths,
		Variable:    variable,
	}
}

func GenerateSingleManager(groupIds []string, clients []ClientConnection) (Manager, error) {
	images, err := generateBase64Images("images/profileImages", 1)
	if err != nil {
		printers["warning"].Printf("Unable to generate base64 image in Clients generator %s\n", err.Error())
	}

	managers := make([]map[string]string, len(clients))
	for i := 0; i < len(clients); i++ {
		managers[i] = map[string]string{
			"client": "",
		}
	}

	valueData := map[string]interface{}{
		"name":       faker.FirstName(),
		"lastName":   faker.LastName(),
		"address":    faker.GetRealAddress().Address,
		"city":       faker.GetRealAddress().City,
		"postalCode": faker.GetRealAddress().PostalCode,
		"managers":   managers,
	}

	if len(images) == 1 {
		valueData["profileImage"] = images[0]
	}

	connections := make([]map[string]string, len(clients))
	for i, conn := range clients {
		connections[i] = map[string]string{
			"name":          fmt.Sprintf("managers.%d.client", i),
			"structureType": conn.StructureType,
			"variableId":    conn.VariableID,
		}
	}

	b, err := json.Marshal(valueData)
	if err != nil {
		return Manager{}, err
	}

	uniqueName := uuid.New().String()
	return newManager(
		uniqueName,
		connections,
		nil,
		newManagerVariable(uniqueName, "eng", "modifiable", "", string(b), pickRandomUniqueGroups(groupIds, 3)),
	), nil
}
