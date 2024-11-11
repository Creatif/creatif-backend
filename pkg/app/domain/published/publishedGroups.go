package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type PublishedGroups struct {
	ID        string    `gorm:"primarykey;type:text"`
	VersionID string    `gorm:"primarykey:type:text"`
	ProjectID string    `gorm:"primarykey;type:text"`
	Name      string    `gorm:"type:text"`
	FileName  string    `gorm:"type:text"`
	MimeType  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func (PublishedGroups) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_GROUPS_TABLE)
}
