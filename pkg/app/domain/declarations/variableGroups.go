package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
)

type VariableGroup struct {
	GroupID    string         `gorm:"type:text"`
	VariableID string         `gorm:"type:text"`
	Groups     pq.StringArray `gorm:"type:text[];not_null"`
}

func NewVariableGroup(groupId, variableId string, groups pq.StringArray) VariableGroup {
	return VariableGroup{
		GroupID:    groupId,
		VariableID: variableId,
		Groups:     groups,
	}
}

func (VariableGroup) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.VARIABLE_GROUPS_TABLE)
}
