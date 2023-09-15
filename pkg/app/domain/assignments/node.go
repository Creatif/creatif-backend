package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// text, image, file, json, code, boolean
type Node struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name string `gorm:"index;uniqueIndex:unique_node"`

	// TODO: change this to be string when projects and exploration are over, project must exist and be UUID
	// TODO: remove pointer since now it allows null
	/*	ProjectID *string `gorm:"uniqueIndex:unique_node"`
		Project   domain.Project*/

	DeclarationNodeID string            `gorm:"type:text CHECK(length(declaration_node_id)=26)"`
	DeclarationNode   declarations.Node `gorm:"foreignKey:DeclarationNodeID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"<-:createNode"`
	UpdatedAt time.Time
}

func NewNode(name string, declarationNodeID string) Node {
	return Node{
		Name:              name,
		DeclarationNodeID: declarationNodeID,
	}
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (Node) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODES_TABLE)
}
