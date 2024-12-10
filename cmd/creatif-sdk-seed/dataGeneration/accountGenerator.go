package dataGeneration

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"os"
)

var printers = map[string]*color.Color{
	"warning": color.New(color.FgYellow).Add(color.Bold),
}

type AccountVariable struct {
	Name      string
	Locale    string
	Behaviour string
	Groups    []string
	Metadata  string
	Value     string
}

type Account struct {
	Name        string
	Connections []map[string]string
	ImagePaths  []string
	Variable    AccountVariable
}

func newAccountVariable(name, locale, behaviour, metadata, value string, groups []string) AccountVariable {
	return AccountVariable{
		Name:      name,
		Locale:    locale,
		Behaviour: behaviour,
		Groups:    groups,
		Metadata:  metadata,
		Value:     value,
	}
}

func generateBase64Images(dir string, numOfImages int) ([]string, error) {
	images, err := pickImages(dir, numOfImages)
	if err != nil {
		return nil, err
	}

	b64Images := make([]string, numOfImages)
	for i, img := range images {
		imagePath := fmt.Sprintf("%s/%s", dir, img.Name())
		bytes, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, err
		}

		b64 := fmt.Sprintf("data:image/webp#webp;base64,%s", base64.StdEncoding.EncodeToString(bytes))

		b64Images[i] = b64
	}

	return b64Images, nil
}

func GenerateSingleAccount(groupIds []string) (Account, error) {
	images, err := generateBase64Images("images/profileImages", 1)
	if err != nil {
		printers["warning"].Printf("Unable to generate base64 image in Accounts generator %s\n", err.Error())
	}

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
		return Account{}, err
	}

	uniqueName := uuid.New().String()
	return newAccount(
		uniqueName,
		nil,
		[]string{"profileImage"},
		newAccountVariable(uniqueName, "eng", "modifiable", "", string(b), pickRandomUniqueGroups(groupIds, 3)),
	), nil
}

func newAccount(name string, connections []map[string]string, imagePaths []string, variable AccountVariable) Account {
	return Account{
		Name:        name,
		Connections: connections,
		ImagePaths:  imagePaths,
		Variable:    variable,
	}
}
