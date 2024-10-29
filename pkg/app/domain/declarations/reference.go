package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"time"
)

type Reference struct {
	ID        string `gorm:"primarykey;type:text"`
	ProjectID string `gorm:"index"`

	Name       string `gorm:"type:text"`
	ParentType string `gorm:"type:text"`
	ChildType  string `gorm:"type:text"`

	ParentStructureID string `gorm:"index;type:text"`
	ChildStructureID  string `gorm:"index;type:text"`

	// must be structure type item
	ParentID string `gorm:"index;type:text"`
	// must be entire structure
	ChildID string `gorm:"index;type:text"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewReference(name, parentType, childType, parentId, childId, parentStructureId, childStructureId, projectId string) Reference {
	return Reference{
		ID:                ksuid.New().String(),
		Name:              name,
		ParentType:        parentType,
		ParentStructureID: parentStructureId,
		ProjectID:         projectId,
		ChildStructureID:  childStructureId,
		ChildType:         childType,
		ParentID:          parentId,
		ChildID:           childId,
	}
}

func (Reference) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.REFERENCE_TABLES)
}
