package main

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/lists/createList"
	addToMap2 "creatif/pkg/app/services/maps/addToMap"
	"creatif/pkg/app/services/maps/mapCreate"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
)

const apiKey = "$2a$10$aUlSZKvCLkbA65wWB5tme.a6nQDwJRzJrjm.DAlpD9/m4hjcrgf/u"
const projectId = "01HP45ME16HK3SRC735Q0KWE06"

func main() {
	loadEnv()
	runDb()
	runLogger()

	if err := releaseAllLocks(); err != nil {
		sqlDB, err := storage.SQLDB()
		if err != nil {
			log.Fatalln("Unable to get storage.SQLDB()", err)
		}

		if err := sqlDB.Close(); err != nil {
			log.Fatalln("Unable to disconnect from the database", err)
		}
	}

	if err := loadLocales(); err != nil {
		sqlDB, err := storage.SQLDB()
		if err != nil {
			log.Fatalln("Unable to get storage.SQLDB()", err)
		}

		if err := sqlDB.Close(); err != nil {
			log.Fatalln("Unable to disconnect from the database", err)
		}
	}

	seed()
}

func seed() {
	createGroups(projectId)
	structureNames := map[string][]string{
		"list": []string{"Languages"},
		"map":  []string{"Decks"},
	}

	fmt.Println("Creating structures...")
	mapStructures := make([]mapCreate.View, 0)
	listStructures := make([]createList.View, 0)
	for key, value := range structureNames {
		if key == "list" {
			for _, name := range value {
				listStructures = append(listStructures, listCreate(name))
			}
		}

		if key == "map" {
			for _, name := range value {
				mapStructures = append(mapStructures, createMap(name))
			}
		}
	}

	fmt.Println("Structures finished!")

	englishId := listAdd(listStructures[0].ID, "English", []shared.Reference{})
	frenchId := listAdd(listStructures[0].ID, "French", []shared.Reference{})

	fmt.Println("Creating languages...")
	m := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		m.Add(1)
		go func() {
			for a := 0; a < 1000; a++ {
				languageId := englishId
				if a%2 == 0 {
					languageId = frenchId
				}

				addToMap(mapStructures[0].ID, uuid.NewString(), []shared.Reference{
					{
						Name:          "language",
						StructureType: "list",
						StructureName: "Languages",
						VariableID:    languageId,
					},
				})
			}

			m.Done()
		}()
	}
	m.Wait()
}

func createGroups(projectId string) {
	l := logger.NewLogBuilder()
	groups := make([]string, 0)
	for i := 0; i < 100; i++ {
		groups = append(groups, fmt.Sprintf("group-%d", i))
	}

	handler := addGroups.New(addGroups.NewModel(projectId, groups), auth.NewNoopAuthentication(), l)

	_, err := handler.Handle()
	if err != nil {
		log.Fatalln(err)
	}
}

func createMap(name string) mapCreate.View {
	l := logger.NewLogBuilder()
	handler := mapCreate.New(mapCreate.NewModel(projectId, name, []mapCreate.VariableModel{}), auth.NewNoopAuthentication(), l)
	m, err := handler.Handle()
	if err != nil {
		log.Fatalln(err)
	}

	return m
}

func listCreate(name string) createList.View {
	l := logger.NewLogBuilder()
	handler := createList.New(createList.NewModel(projectId, name, []createList.Variable{}), auth.NewNoopAuthentication(), l)
	m, err := handler.Handle()
	if err != nil {
		log.Fatalln(err)
	}

	return m
}

func addToMap(structureId, variableName string, references []shared.Reference) string {
	l := logger.NewLogBuilder()
	value := map[string]interface{}{
		"name": variableName,
	}

	v, err := json.Marshal(value)
	if err != nil {
		log.Fatalln(err)
	}

	handler := addToMap2.New(addToMap2.NewModel(projectId, structureId, addToMap2.VariableModel{
		Name:      variableName,
		Metadata:  nil,
		Locale:    "eng",
		Groups:    []string{},
		Behaviour: "modifiable",
		Value:     v,
	}, references), auth.NewNoopAuthentication(), l)

	entry, err := handler.Handle()
	if err != nil {
		validationError, ok := err.(appErrors.AppError[map[string]string])
		if ok {
			log.Fatalln(validationError.Data())
		}
		log.Fatalln(err)
	}

	return entry.ID
}

func listAdd(structureId, variableName string, references []shared.Reference) string {
	l := logger.NewLogBuilder()
	value := map[string]interface{}{
		"name": variableName,
	}

	v, err := json.Marshal(value)
	if err != nil {
		log.Fatalln(err)
	}

	handler := addToList.New(addToList.NewModel(projectId, structureId, addToList.VariableModel{
		Name:      variableName,
		Metadata:  nil,
		Locale:    "eng",
		Groups:    []string{},
		Behaviour: "modifiable",
		Value:     v,
	}, references), auth.NewNoopAuthentication(), l)

	entry, err := handler.Handle()
	if err != nil {
		validationError, ok := err.(appErrors.AppError[map[string]string])
		if ok {
			log.Fatalln(validationError.Data())
		}
		log.Fatalln(err)
	}

	return entry.ID
}
