package declarations

import "github.com/microcosm-cc/bluemonday"

type CreateNode struct {
	Name       string                    `json:"name"`
	Type       string                    `json:"type"`
	Groups     []string                  `json:"groups"`
	Behaviour  string                    `json:"behaviour"`
	Validation map[string]NodeValidation `json:"validation"`
	Metadata   string                    `json:"metadata"`
}

type NodeValidation struct {
	Required    bool
	Length      string
	ExactValue  string
	ExactValues string
	IsDate      bool
}

func SanitizeNode(model CreateNode) CreateNode {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Type = p.Sanitize(model.Type)
	model.Behaviour = p.Sanitize(model.Behaviour)

	return model
}
