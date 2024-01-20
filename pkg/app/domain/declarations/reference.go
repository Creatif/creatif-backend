package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Reference struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

	Name       string `gorm:"type:text"`
	ParentType string `gorm:"type:text"`
	ChildType  string `gorm:"type:text"`

	// must be structure type item
	ParentID string `gorm:"type:text"`
	// must be entire structure
	ChildID string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func NewReference(name, parentType, childType, parentId, childId string) Reference {
	return Reference{
		Name:       name,
		ParentType: parentType,
		ChildType:  childType,
		ParentID:   parentId,
		ChildID:    childId,
	}
}

func (Reference) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.REFERENCE_TABLES)
}
