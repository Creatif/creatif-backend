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

	projectids := []string{
		"01HQJ4S0VRSC169QDDY8DJY90H",
	}

	for _, p := range projectids {
		seed(p)
	}
}

func seed(projectId string) {
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
				listStructures = append(listStructures, listCreate(projectId, name))
			}
		}

		if key == "map" {
			for _, name := range value {
				mapStructures = append(mapStructures, createMap(projectId, name))
			}
		}
	}

	fmt.Println("Structures finished!")

	englishId := listAdd(projectId, listStructures[0].ID, "English", []shared.Reference{})
	frenchId := listAdd(projectId, listStructures[0].ID, "French", []shared.Reference{})

	fmt.Println("Creating languages...")
	for i := 0; i < 50; i++ {
		fmt.Println(fmt.Sprintf("Batch %d finished.", i))
		addBatch(projectId, englishId, frenchId, mapStructures[0].ID)
	}
}

func addBatch(projectId, englishId, frenchId, mapStructureId string) {
	m := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		m.Add(1)
		go func() {
			for a := 0; a < 5000; a++ {
				languageId := englishId
				if a%2 == 0 {
					languageId = frenchId
				}

				addToMap(projectId, mapStructureId, uuid.NewString(), []shared.Reference{
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
	groups := make([]addGroups.GroupModel, 0)
	for i := 0; i < 100; i++ {
		groups = append(groups, addGroups.GroupModel{
			ID:     "",
			Name:   fmt.Sprintf("group-%d", i),
			Type:   "new",
			Action: "create",
		})
	}

	handler := addGroups.New(addGroups.NewModel(projectId, groups), auth.NewNoopAuthentication(), l)

	_, err := handler.Handle()
	l.Flush("")
	if err != nil {
		log.Fatalln(err)
	}
}

func createMap(projectId, name string) mapCreate.View {
	l := logger.NewLogBuilder()
	handler := mapCreate.New(mapCreate.NewModel(projectId, name, []mapCreate.VariableModel{}), auth.NewNoopAuthentication(), l)
	m, err := handler.Handle()
	if err != nil {
		log.Fatalln(err)
	}
	l.Flush("")

	return m
}

func listCreate(projectId, name string) createList.View {
	l := logger.NewLogBuilder()
	handler := createList.New(createList.NewModel(projectId, name, []createList.Variable{}), auth.NewNoopAuthentication(), l)
	m, err := handler.Handle()
	if err != nil {
		log.Fatalln(err)
	}
	l.Flush("")

	return m
}

func addToMap(projectId, structureId, variableName string, references []shared.Reference) string {
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
	l.Flush("")

	return entry.ID
}

func listAdd(projectId, structureId, variableName string, references []shared.Reference) string {
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
	l.Flush("")

	return entry.ID
}
