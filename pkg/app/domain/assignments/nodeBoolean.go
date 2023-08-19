package assignments

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NodeBoolean struct {
	ID               string `gorm:"primarykey"`
	AssignmentNodeID string `gorm:"primarykey;autoincrement:false"`

	Value bool
}

func NewNodeBoolean(assignmentNodeID string, value bool) NodeBoolean {
	return NodeBoolean{
		Value:            value,
		AssignmentNodeID: assignmentNodeID,
	}
}

func (u *NodeBoolean) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (NodeBoolean) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODE_BOOLEAN_TABLE)
}
