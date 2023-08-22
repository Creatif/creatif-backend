package assignments

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ValueNode struct {
	ID               string `gorm:"primarykey"`
	AssignmentNodeID string `gorm:"primarykey;autoincrement:false"`

	Value datatypes.JSON
}

func NewValueNode(assignmentNodeID string, value datatypes.JSON) ValueNode {
	return ValueNode{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *ValueNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (ValueNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_VALUE_NODE)
}
