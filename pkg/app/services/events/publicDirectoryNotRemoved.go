package events

import (
	"creatif/pkg/lib/constants"
	"encoding/json"
)

type PublicDirectoryNotRemoved struct {
	FilePath   string `json:"filePath"`
	RelationID string `json:"relationID"`
	ProjectID  string `json:"projectID"`
}

func NewPublicDirectoryNotRemoved(filePath, relation, projectID string) FileNotRemovedEvent {
	return FileNotRemovedEvent{
		FilePath:   filePath,
		RelationID: relation,
		ProjectID:  projectID,
	}
}

func (e PublicDirectoryNotRemoved) Project() string {
	return e.ProjectID
}

func (e PublicDirectoryNotRemoved) Type() string {
	return constants.FileNotRemovedEvent
}

func (e PublicDirectoryNotRemoved) Data() []byte {
	b, _ := json.Marshal(e)
	return b
}
