package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type CreateMap struct {
	Nodes []string `json:"nodes"`
	Name  string   `json:"name"`
}

func SanitizeMapModel(model CreateMap) CreateMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	newNodes := make([]string, 0)
	for _, n := range model.Nodes {
		newNodes = append(newNodes, p.Sanitize(n))
	}
	model.Nodes = newNodes

	return model
}
