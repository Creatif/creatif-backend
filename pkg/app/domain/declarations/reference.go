package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Reference struct {
	ID string `gorm:"primarykey;type:text;default:gen_ulid()"`

	ParentType string `gorm:"type:text"`
	ChildType  string `gorm:"type:text"`

	// must be structure type item
	ParentID      string `gorm:"type:text"`
	ParentShortID string `gorm:"type:text"`
	// must be entire structure
	ChildID      string `gorm:"type:text"`
	ChildShortID string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func NewReference(parentType, childType, parentId, parentShortId, childId, childShortId string) Reference {
	return Reference{
		ParentType:    parentType,
		ChildType:     childType,
		ParentID:      parentId,
		ParentShortID: parentShortId,
		ChildID:       childId,
		ChildShortID:  childShortId,
	}
}

func (Reference) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.REFERENCE_TABLES)
}
