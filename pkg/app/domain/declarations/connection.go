package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Connection struct {
	ProjectID string `gorm:"index"`

	Path                string `gorm:"type:text"`
	ParentVariableID    string `gorm:"type:text"`
	ParentStructureType string `gorm:"type:text"`

	ChildVariableID    string `gorm:"index;type:text"`
	ChildStructureType string `gorm:"index;type:text"`

	CreatedAt time.Time `gorm:"<-:create"`
}

func NewConnection(projectId, path, parentVariableId, parentStructureType, childVariableId, childStructureType string) Connection {
	return Connection{
		ProjectID:           projectId,
		Path:                path,
		ParentVariableID:    parentVariableId,
		ParentStructureType: parentStructureType,
		ChildVariableID:     childVariableId,
		ChildStructureType:  childStructureType,
	}
}

func (Connection) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.CONNECTIONS_TABLE)
}
