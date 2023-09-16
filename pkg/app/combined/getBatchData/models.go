package getBatchData

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"time"
)

type variable struct {
	Name string
	Type string
}

type Model struct {
	Variables []variable

	mapIds      []string
	variableIds []string
}

func NewModel(variables map[string]string) *Model {
	models := make([]variable, len(variables))
	count := 0
	for name, t := range variables {
		models[count] = variable{
			Name: name,
			Type: t,
		}
		count++
	}

	return &Model{
		Variables: models,
	}
}

type View struct {
	ID string `json:"id"`

	Name      string         `json:"name"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups"`
	Metadata  interface{}    `json:"metadata"`
	Value     interface{}    `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model map[string]interface{}) map[string]interface{} {
	view := make(map[string]interface{})
	variables := model["variables"]
	convertedVariables, ok := variables.([]Variable)
	if ok {
		variableView := make(map[string][]View)
		for _, n := range convertedVariables {
			if _, ok := variableView[n.Name]; !ok {
				variableView[n.Name] = make([]View, 0)
			}

			variableView[n.Name] = append(variableView[n.Name], View{
				ID:        n.ID,
				Name:      n.Name,
				Behaviour: n.Behaviour,
				Groups:    n.Groups,
				Metadata:  n.Metadata,
				Value:     n.Value,
				CreatedAt: n.CreatedAt,
				UpdatedAt: n.UpdatedAt,
			})
		}

		view["variables"] = variableView
	}

	maps := model["maps"]
	convertedMaps, ok := maps.(map[string][]Variable)

	if ok {
		resolvedMaps := make(map[string][]View)
		for key, mapVariables := range convertedMaps {
			a := make([]View, 0)
			for _, n := range mapVariables {
				a = append(a, View{
					ID:        n.ID,
					Name:      n.Name,
					Behaviour: n.Behaviour,
					Groups:    n.Groups,
					Metadata:  n.Metadata,
					Value:     n.Value,
					CreatedAt: n.CreatedAt,
					UpdatedAt: n.UpdatedAt,
				})
			}

			resolvedMaps[key] = a
		}

		view["maps"] = resolvedMaps
	}

	return view
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"validVariables": a.Variables,
		"validNames":     a.Variables,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("validVariables", validation.By(func(value interface{}) error {
				variables := value.([]variable)

				for _, t := range variables {
					if t.Name == "" {
						return errors.New("Variable name cannot be empty")
					}
				}

				for _, t := range variables {
					if t.Type != "variable" && t.Type != "map" {
						return errors.New(fmt.Sprintf("Invalid type in variable with name '%s'. Valid type are 'map' and 'variable'", t.Name))
					}
				}

				return nil
			})),
			validation.Key("validNames", validation.By(func(value interface{}) error {
				variables := value.([]variable)

				variablesNames := sdk.Filter(sdk.Map(variables, func(idx int, value variable) string {
					if value.Type == "variable" {
						return value.Name
					}

					return ""
				}), func(idx int, value string) bool {
					return value != ""
				})

				mapNames := sdk.Filter(sdk.Map(variables, func(idx int, value variable) string {
					if value.Type == "map" {
						return value.Name
					}

					return ""
				}), func(idx int, value string) bool {
					return value != ""
				})

				var foundVariables []declarations.Variable
				if res := storage.Gorm().Table((declarations.Variable{}).TableName()).Select("ID").Where("name IN (?)", variablesNames).Find(&foundVariables); res.Error != nil {
					return errors.New("One of the variables or map names given is invalid or does not exist.")
				}

				var maps []declarations.Map
				if res := storage.Gorm().Table((declarations.Map{}).TableName()).Select("ID").Where("name IN (?)", mapNames).Find(&maps); res.Error != nil {
					return errors.New("One of the variables or map names given is invalid or does not exist.")
				}

				if (len(variablesNames) + len(mapNames)) != len(variables) {
					return errors.New("One of the variables or map names given is invalid or does not exist.")
				}

				a.variableIds = sdk.Map(foundVariables, func(idx int, value declarations.Variable) string {
					return value.ID
				})

				a.mapIds = sdk.Map(maps, func(idx int, value declarations.Map) string {
					return value.ID
				})

				return nil
			})),
		),
	); err != nil {
		var e map[string]string
		b, err := json.Marshal(err)
		if err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		if err := json.Unmarshal(b, &e); err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		return e
	}

	return nil
}
