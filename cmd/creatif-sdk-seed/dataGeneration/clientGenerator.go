package dataGeneration

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

var printers = map[string]*color.Color{
	"warning": color.New(color.FgYellow).Add(color.Bold),
}

type ClientVariable struct {
	Name      string
	Locale    string
	Behaviour string
	Groups    []string
	Metadata  string
	Value     string
}

type Client struct {
	Name        string
	Connections []map[string]string
	ImagePaths  []string
	Variable    ClientVariable
}

func newClientVariable(name, locale, behaviour, metadata, value string, groups []string) ClientVariable {
	return ClientVariable{
		Name:      name,
		Locale:    locale,
		Behaviour: behaviour,
		Groups:    groups,
		Metadata:  metadata,
		Value:     value,
	}
}

func newClient(name string, connections []map[string]string, imagePaths []string, variable ClientVariable) Client {
	return Client{
		Name:        name,
		Connections: connections,
		ImagePaths:  imagePaths,
		Variable:    variable,
	}
}

func GenerateSingleClient(groupIds []string) (Client, error) {
	images, err := generateBase64Images("images/profileImages", 1)
	if err != nil {
		printers["warning"].Printf("Unable to generate base64 image in Clients generator %s\n", err.Error())
	}

	clientValueData := map[string]string{
		"name":       faker.FirstName(),
		"lastName":   faker.LastName(),
		"address":    faker.GetRealAddress().Address,
		"city":       faker.GetRealAddress().City,
		"postalCode": faker.GetRealAddress().PostalCode,
	}

	if len(images) == 1 {
		clientValueData["profileImage"] = images[0]
	}

	b, err := json.Marshal(clientValueData)
	if err != nil {
		return Client{}, err
	}

	uniqueName := uuid.New().String()
	return newClient(
		uniqueName,
		nil,
		[]string{"profileImage"},
		newClientVariable(uniqueName, "eng", "modifiable", "", string(b), pickRandomUniqueGroups(groupIds, 3)),
	), nil
}
