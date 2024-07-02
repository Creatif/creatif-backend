package events

import (
	"creatif/pkg/lib/constants"
	"encoding/json"
)

type FileNotRemovedEvent struct {
	FilePath   string `json:"filePath"`
	RelationID string `json:"relationID"`
}

func NewFileNotRemoveEvent(filePath, relation string) FileNotRemovedEvent {
	return FileNotRemovedEvent{
		FilePath:   filePath,
		RelationID: relation,
	}
}

func (e FileNotRemovedEvent) Type() string {
	return constants.FileNotRemovedEvent
}

func (e FileNotRemovedEvent) Data() []byte {
	b, _ := json.Marshal(e)
	return b
}
