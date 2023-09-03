package assignments

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/app/domain/declarations"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

// text, image, file, json, code, boolean
type Node struct {
	ID ksuid.KSUID `gorm:"primarykey;type:text CHECK(length(id)=27)"`

	Name string `gorm:"index;uniqueIndex:unique_node"`

	// TODO: change this to be ksuid.KSUID when projects and exploration are over, project must exist and be UUID
	// TODO: remove pointer since now it allows null
	/*	ProjectID *string `gorm:"uniqueIndex:unique_node"`
		Project   domain.Project*/

	DeclarationNodeID ksuid.KSUID       `gorm:"type:text CHECK(length(declaration_node_id)=27)"`
	DeclarationNode   declarations.Node `gorm:"foreignKey:DeclarationNodeID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewNode(name string, declarationNodeID ksuid.KSUID) Node {
	return Node{
		Name:              name,
		DeclarationNodeID: declarationNodeID,
	}
}

func (u *Node) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = ksuid.New()

	return
}

func (Node) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_NODES_TABLE)
}
