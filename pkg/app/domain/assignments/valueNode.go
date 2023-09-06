package assignments

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ValueNode struct {
	ID               ksuid.KSUID `gorm:"primarykey;type:text CHECK(length(id)=27)"`
	AssignmentNodeID ksuid.KSUID `gorm:"type:text CHECK(length(assignment_node_id)=27);constraint:OnDelete:CASCADE"`

	Value datatypes.JSON
}

func NewValueNode(assignmentNodeID ksuid.KSUID, value datatypes.JSON) ValueNode {
	return ValueNode{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *ValueNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = ksuid.New()

	return
}

func (ValueNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_VALUE_NODE)
}
