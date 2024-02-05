package main

import (
	"creatif/pkg/app/auth"
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
const projectId = "01HNWQ28QYABWEHNP3R9BH6YXV"

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
	structureNames := []string{
		"Languages",
		"Decks",
	}

	fmt.Println("Creating structures...")
	structures := make([]mapCreate.View, len(structureNames))
	for i, s := range structureNames {
		structures[i] = createMap(s)
	}
	fmt.Println("Structures finished!")

	englishId := addToMap(structures[0].ID, "English", []shared.Reference{})
	frenchId := addToMap(structures[0].ID, "French", []shared.Reference{})

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
				addToMap(structures[1].ID, uuid.NewString(), []shared.Reference{
					{
						Name:          "language",
						StructureType: "map",
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

func createMap(name string) mapCreate.View {
	l := logger.NewLogBuilder()
	handler := mapCreate.New(mapCreate.NewModel(projectId, name, []mapCreate.VariableModel{}), auth.NewNoopAuthentication(), l)
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
		Groups:    []string{variableName, "default"},
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
