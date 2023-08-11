package request

import "github.com/microcosm-cc/bluemonday"

type CreateNode struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Group     string `json:"group"`
	Behaviour string `json:"behaviour"`
}

func SanitizeNode(model CreateNode) CreateNode {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Type = p.Sanitize(model.Type)
	model.Group = p.Sanitize(model.Group)
	model.Behaviour = p.Sanitize(model.Behaviour)

	return model
}
