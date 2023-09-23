package declarations

import (
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"errors"
	"github.com/microcosm-cc/bluemonday"
)

type MapVariableModel struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Value     string   `json:"value"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
}

type Entry struct {
	Type  string
	Model interface{}
}

type CreateMap struct {
	Entries   []Entry `json:"entries"`
	Name      string  `json:"name"`
	ProjectID string  `param:"projectID"`
}

func (u *CreateMap) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	if _, ok := v["name"].(string); !ok {
		return errors.New("'name' is not a string")
	}

	rawEntries, ok := v["entries"]
	if !ok {
		return errors.New("'entries' are missing")
	}

	var entries []Entry
	if err := sdk.ConvertByUnmarshaling(rawEntries, &entries); err != nil {
		return err
	}

	newEntries := make([]Entry, 0)
	for _, entry := range entries {
		if entry.Type == "variable" {
			var variable MapVariableModel
			if err := sdk.ConvertByUnmarshaling(entry.Model, &variable); err != nil {
				return err
			}

			newEntries = append(newEntries, Entry{
				Type:  entry.Type,
				Model: variable,
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
	model.ProjectID = p.Sanitize(model.ProjectID)

	entries := model.Entries
	newEntries := make([]Entry, len(entries))
	for i := 0; i < len(entries); i++ {
		newEntry := Entry{
			Type:  p.Sanitize(entries[i].Type),
			Model: nil,
		}

		if entries[i].Type == "variable" {
			m := entries[i].Model.(MapVariableModel)
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
