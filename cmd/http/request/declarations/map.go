package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type CreateMap struct {
	Entries []Entry `json:"entries"`
	Name    string  `json:"name"`
}

type MapNodeModel struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
}

type Entry struct {
	Type  string
	Model interface{}
}

func (u *CreateMap) UnmarshalJSON(b []byte) error {
	return nil
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
