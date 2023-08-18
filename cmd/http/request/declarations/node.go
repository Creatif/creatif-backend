package declarations

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type CreateNode struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Groups     []string       `json:"groups"`
	Behaviour  string         `json:"behaviour"`
	Validation NodeValidation `json:"validation"`
	Metadata   string         `json:"metadata"`
}

type ValidationLength struct {
	Min   int `json:"min"`
	Max   int `json:"max"`
	Exact int `json:"exact"`
}

type NodeValidation struct {
	Required    bool             `json:"metadata"`
	Length      ValidationLength `json:"length"`
	ExactValue  string           `json:"exactValue"`
	ExactValues []string         `json:"exactValues"`
	IsDate      bool             `json:"isDate"`
}

func SanitizeNode(model CreateNode) CreateNode {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Type = p.Sanitize(model.Type)
	model.Behaviour = p.Sanitize(model.Behaviour)

	model.Groups = sdk.Sanitize(model.Groups, func(k int, v string) string {
		return p.Sanitize(v)
	})

	newNodeValidation := NodeValidation{}
	newNodeValidation.Required = model.Validation.Required
	newNodeValidation.IsDate = model.Validation.IsDate
	newNodeValidation.ExactValue = p.Sanitize(model.Validation.ExactValue)

	newNodeValidation.ExactValues = sdk.Sanitize(model.Validation.ExactValues, func(k int, v string) string {
		return p.Sanitize(v)
	})

	model.Validation = newNodeValidation

	return model
}
