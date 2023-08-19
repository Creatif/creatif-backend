package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// text, image, file, json, code, boolean
type Node struct {
	ID string `gorm:"primarykey"`

	Name string `gorm:"index;uniqueIndex:unique_node"`

	// TODO: change this to be uuid.UUID when projects and exploration are over, project must exist and be UUID
	// TODO: remove pointer since now it allows null
	/*	ProjectID *string `gorm:"uniqueIndex:unique_node"`
		Project   domain.Project*/

	ValueType string

	DeclarationNodeID string
	DeclarationNode   declarations.Node `gorm:"foreignKey:DeclarationNodeID"`

	NodeText    *NodeText    `gorm:"-"`
	NodeBoolean *NodeBoolean `gorm:"-"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewNode(name string, declarationNodeID string) Node {
	return Node{
		Name:              name,
		DeclarationNodeID: declarationNodeID,
	}
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}

func (Node) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODES_TABLE)
}
