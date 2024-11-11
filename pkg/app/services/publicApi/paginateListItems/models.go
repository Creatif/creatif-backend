package paginateListItems

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared/queryProcessor"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
	"time"
)

type Model struct {
	ProjectID     string
	StructureName string
	VersionName   string
	Options       Options

	Page    int
	Limit   int
	Order   string
	SortBy  string
	Search  string
	Locales []string
	Groups  []string
	Query   []queryProcessor.Query
}

type Options struct {
	ValueOnly bool
}

func NewModel(versionName, projectId, structureName string, page, limit int, order string, sortBy, search string, lcls, groups []string, options Options, query []queryProcessor.Query) Model {
	return Model{
		StructureName: structureName,
		Options:       options,
		VersionName:   versionName,
		ProjectID:     projectId,
		Page:          page,
		Limit:         limit,
		Order:         order,
		SortBy:        sortBy,
		Search:        search,
		Locales:       lcls,
		Groups:        groups,
		Query:         query,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":     a.ProjectID,
		"order":         a.Order,
		"sortBy":        a.SortBy,
		"structureName": a.StructureName,
		"locales":       a.Locales,
		"versionName":   a.VersionName,
		"limit":         a.Limit,
		"query":         a.Query,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("versionName", validation.When(a.VersionName != "", validation.RuneLength(1, 200))),
			validation.Key("structureName", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("order", validation.By(func(value interface{}) error {
				order := strings.ToLower(a.Order)
				if order != "desc" && order != "asc" {
					return errors.New("Order must be either DESC or ASC")
				}
				return nil
			})),
			validation.Key("sortBy", validation.By(func(value interface{}) error {
				sortBy := strings.ToLower(a.SortBy)
				if sortBy != "name" && sortBy != "created_at" && sortBy != "updated_at" && sortBy != "index" {
					return errors.New("Invalid sortBy field. sortBy can be: name, created_at, updated_at, index")
				}
				return nil
			})),
			validation.Key("locales", validation.By(func(value interface{}) error {
				for _, l := range a.Locales {
					_, err := locales.GetIDWithAlpha(l)
					if err != nil {
						return errors.New(fmt.Sprintf("Locale %s does not exist.", l))
					}
				}
				return nil
			})),
			validation.Key("limit", validation.By(func(value interface{}) error {
				limit := value.(int)

				if limit > 1000 {
					return errors.New("Limit cannot be bigger than 1000")
				}

				return nil
			})),
			validation.Key("query", validation.By(func(value interface{}) error {
				query := value.([]queryProcessor.Query)

				for _, q := range query {
					if q.Column == "" {
						return errors.New("Query 'column' cannot be empty")
					}

					if q.Operator == "" {
						return errors.New("Query 'operator' cannot be empty")
					}

					validateOperators := func() error {
						validOperators := []string{"equal", "unequal", "greaterThan", "lessThan", "greaterThanOrEqual", "lessThanOrEqual"}
						found := false
						for _, v := range validOperators {
							if v == q.Operator {
								found = true
								break
							}
						}

						if !found {
							return errors.New(fmt.Sprintf("Invalid operator. Operators can be only %s", strings.Join(validOperators, ", ")))
						}

						return nil
					}

					validateTypes := func() error {
						validTypes := []string{"int", "float", "string"}
						found := false
						for _, v := range validTypes {
							if v == q.Type {
								found = true
								break
							}
						}

						if !found {
							return errors.New(fmt.Sprintf("Invalid data type. Data type can be only %s", strings.Join(validTypes, ", ")))
						}

						return nil
					}

					if err := validateOperators(); err != nil {
						return err
					}

					if err := validateTypes(); err != nil {
						return err
					}
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	StructureID      string `json:"structureId"`
	StructureShortID string `json:"structureShortId"`
	StructureName    string `json:"structureName"`

	Name    string `json:"itemName"`
	ID      string `json:"itemId"`
	ShortID string `json:"itemShortId"`

	ProjectID string      `json:"projectId"`
	Locale    string      `json:"locale"`
	Index     float64     `json:"index"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LogicModel struct {
	Items   []Item
	Options Options
}

func newView(model LogicModel) interface{} {
	if model.Options.ValueOnly {
		returnValue := make([]map[string]interface{}, len(model.Items))
		for i, val := range model.Items {
			var m map[string]interface{}
			// ok to ignore
			json.Unmarshal(val.Value, &m)

			returnValue[i] = m
		}

		return returnValue
	}

	views := make([]View, len(model.Items))
	for i, item := range model.Items {
		locale, _ := locales.GetAlphaWithID(item.Locale)

		views[i] = View{
			StructureID:      item.ID,
			StructureShortID: item.ShortID,
			StructureName:    item.StructureName,
			Name:             item.ItemName,
			ID:               item.ItemID,
			ShortID:          item.ItemShortID,
			ProjectID:        item.ProjectID,
			Locale:           locale,
			Index:            item.Index,
			Groups:           item.Groups,
			Behaviour:        item.Behaviour,
			Value:            item.Value,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		}
	}

	return views
}
