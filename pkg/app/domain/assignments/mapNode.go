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
type MapNode struct {
	ID uuid.UUID `gorm:"primarykey;type:uuid"`

	Name string `gorm:"index;uniqueIndex:unique_node"`

	// TODO: change this to be uuid.UUID when projects and exploration are over, project must exist and be UUID
	// TODO: remove pointer since now it allows null
	/*	ProjectID *string `gorm:"uniqueIndex:unique_node"`
		Project   domain.Project*/

	MapNodeID uuid.UUID        `gorm:"type:uuid"`
	MapNode   declarations.Map `gorm:"foreignKey:MapNodeID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewMapNode(name string, mapNodeID uuid.UUID) MapNode {
	return MapNode{
		Name:      name,
		MapNodeID: mapNodeID,
	}
}

func (u *MapNode) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (MapNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_MAP_NODES_TABLE)
}
