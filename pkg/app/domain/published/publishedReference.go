package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type PublishedReference struct {
	ID        string `gorm:"primaryKey;type:text"`
	ProjectID string `gorm:"index"`

	VersionID string  `gorm:"primaryKey;type:text"`
	Version   Version `gorm:"foreignKey:VersionID"`

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
	UpdatedAt time.Time `gorm:"<-:update"`
}

func NewPublishedReference(
	id,
	projectId,
	versionId,
	name,
	parentType,
	childType,
	parentStructureId,
	childStructureId,
	parentId,
	childId string,
) PublishedReference {
	return PublishedReference{
		ID:                id,
		ProjectID:         projectId,
		VersionID:         versionId,
		Name:              name,
		ParentType:        parentType,
		ChildType:         childType,
		ParentStructureID: parentStructureId,
		ChildStructureID:  childStructureId,
		ParentID:          parentId,
		ChildID:           childId,
	}
}

func (PublishedReference) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_REFERENCES_TABLE)
}
