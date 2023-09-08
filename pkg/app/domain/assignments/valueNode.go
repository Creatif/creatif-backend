package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ValueNode struct {
	ID               string `gorm:"primarykey;type:text CHECK(length(id)=26)"`
	AssignmentNodeID string `gorm:"type:text CHECK(length(assignment_node_id)=26);constraint:OnDelete:CASCADE"`

	Value datatypes.JSON
}

func NewValueNode(assignmentNodeID string, value datatypes.JSON) ValueNode {
	return ValueNode{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *ValueNode) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (ValueNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_VALUE_NODE)
}
