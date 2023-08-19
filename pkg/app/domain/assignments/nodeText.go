package assignments

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type NodeText struct {
	ID               string `gorm:"primarykey"`
	AssignmentNodeID string `gorm:"primarykey;autoincrement:false"`

	Value datatypes.JSON
}

func NewNodeText(assignmentNodeID string, value []byte) NodeText {
	return NodeText{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *NodeText) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (NodeText) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODE_TEXT_TABLE)
}
