package events

import (
	"creatif/pkg/lib/constants"
	"encoding/json"
)

type FileNotRemovedEvent struct {
	FilePath   string `json:"filePath"`
	RelationID string `json:"relationID"`
	ProjectID  string `json:"projectID"`
}

func NewFileNotRemoveEvent(filePath, relation, projectID string) FileNotRemovedEvent {
	return FileNotRemovedEvent{
		FilePath:   filePath,
		RelationID: relation,
		ProjectID:  projectID,
	}
}

func (e FileNotRemovedEvent) Project() string {
	return e.ProjectID
}

func (e FileNotRemovedEvent) Type() string {
	return constants.FileNotRemovedEvent
}

func (e FileNotRemovedEvent) Data() []byte {
	b, _ := json.Marshal(e)
	return b
}
