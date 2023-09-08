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
type MapNode struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=26)"`

	Name string `gorm:"index;uniqueIndex:unique_node"`

	// TODO: change this to be string when projects and exploration are over, project must exist and be UUID
	// TODO: remove pointer since now it allows null
	/*	ProjectID *string `gorm:"uniqueIndex:unique_node"`
		Project   domain.Project*/

	MapNodeID string           `gorm:"type:uuid"`
	MapNode   declarations.Map `gorm:"foreignKey:MapNodeID"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewMapNode(name string, mapNodeID string) MapNode {
	return MapNode{
		Name:      name,
		MapNodeID: mapNodeID,
	}
}

func (u *MapNode) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sdk.NewULID()
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

func (MapNode) TableName() string {
	return fmt.Sprintf("%s.%s", "assignments", domain.ASSIGNMENT_MAP_NODES_TABLE)
}
