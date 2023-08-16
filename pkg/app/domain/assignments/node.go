package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// text, image, file, json, code, boolean
type Node struct {
	ID string `gorm:"primarykey"`

	Name string `gorm:"index;uniqueIndex:unique_node"`

	// TODO: change this to be uuid.UUID when projects and exploration are over, project must exist and be UUID
	// TODO: remove pointer since now it allows null
	ProjectID *string `gorm:"uniqueIndex:unique_node"`
	Project   domain.Project

	ValueID   string
	ValueType string

	DeclarationNodeID string
	DeclarationNode   declarations.Node `gorm:"foreignKey:DeclarationNodeID"`

	NodeText    *NodeText    `gorm:"-"`
	NodeBoolean *NodeBoolean `gorm:"-"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewNode(name string) Node {
	return Node{
		Name: name,
	}
}

type NodeText struct {
	ID string `gorm:"primarykey"`

	Value datatypes.JSON
}

func NewNodeText(value []byte) NodeText {
	return NodeText{
		Value: value,
	}
}

type NodeBoolean struct {
	ID string `gorm:"primarykey"`

	Value bool
}

func NewNodeBoolean(value bool) NodeBoolean {
	return NodeBoolean{
		Value: value,
	}
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (u *NodeText) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (u *NodeBoolean) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (Node) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODES_TABLE)
}

func (NodeText) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODE_TEXT_TABLE)
}

func (NodeBoolean) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODE_BOOLEAN_TABLE)
}
