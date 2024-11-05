package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/lib/pq"
)

type VariableGroup struct {
	VariableID string         `gorm:"type:text"`
	Groups     pq.StringArray `gorm:"type:text[];not_null"`
}

func NewVariableGroup(variableId string, groups pq.StringArray) VariableGroup {
	return VariableGroup{
		VariableID: variableId,
		Groups:     groups,
	}
}

func (VariableGroup) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.VARIABLE_GROUPS_TABLE)
}
