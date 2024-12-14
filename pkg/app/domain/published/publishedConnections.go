package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type PublishedConnection struct {
	ProjectID string `gorm:"index"`
	VersionID string `gorm:"index"`

	Path                string `gorm:"type:text"`
	ParentVariableID    string `gorm:"type:text"`
	ParentStructureType string `gorm:"type:text"`

	ChildVariableID    string `gorm:"index;type:text"`
	ChildStructureType string `gorm:"index;type:text"`

	CreatedAt time.Time `gorm:"<-:create"`
}

func (PublishedConnection) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_CONNECTIONS_TABLE)
}
