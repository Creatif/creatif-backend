package declarations

import (
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"errors"
	"github.com/microcosm-cc/bluemonday"
)

type MapNodeModel struct {
	Name      string   `json:"name"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
}

type Entry struct {
	Type  string
	Model interface{}
}

type CreateMap struct {
	Entries []Entry `json:"entries"`
	Name    string  `json:"name"`
}

func (u *CreateMap) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}

	if _, ok := v["name"].(string); !ok {
		return errors.New("'name' is not a string")
	}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	rawEntries, ok := v["entries"]
	if !ok {
		return errors.New("'entries' are missing")
	}

	var entries []Entry
	if err := sdk.ConvertByUnmarshaling(rawEntries, entries); err != nil {
		return err
	}

	newEntries := make([]Entry, 0)
	for _, entry := range entries {
		if entry.Type == "node" {
			var node MapNodeModel
			if err := sdk.ConvertByUnmarshaling(entry.Model, node); err != nil {
				return err
			}

			newEntries = append(newEntries, Entry{
				Type:  entry.Type,
				Model: node,
			})
		}
	}

	u.Entries = newEntries
	u.Name = v["name"].(string)

	return nil
}

func SanitizeMapModel(model CreateMap) CreateMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)

	entries := model.Entries
	newEntries := make([]Entry, len(entries))
	for i := 0; i < len(entries); i++ {
		newEntry := Entry{
			Type:  p.Sanitize(entries[i].Type),
			Model: nil,
		}

		if entries[i].Type == "node" {
			m := entries[i].Model.(MapNodeModel)
			m.Name = p.Sanitize(m.Name)
			m.Behaviour = p.Sanitize(m.Behaviour)

			newGroups := make([]string, len(m.Groups))
			for a := 0; a < len(m.Groups); a++ {
				newGroups[a] = p.Sanitize(m.Groups[a])
			}

			m.Groups = newGroups

			newEntry.Model = m
		}

		newEntries = append(newEntries, newEntry)
	}

	model.Entries = newEntries

	return model
}
