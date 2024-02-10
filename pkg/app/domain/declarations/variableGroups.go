package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
)

type VariableGroup struct {
	GroupID    string `gorm:"type:text"`
	VariableID string `gorm:"type:text"`
}

func NewVariableGroup(groupId, variableId string) VariableGroup {
	return VariableGroup{
		GroupID:    groupId,
		VariableID: variableId,
	}
}

func (VariableGroup) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.VARIABLE_GROUPS_TABLE)
}
