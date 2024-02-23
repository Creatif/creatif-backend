package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Version struct {
	ID        string `gorm:"primarykey;type:text;default:gen_ulid()"`
	ProjectID string `gorm:"index;type:text"`

	Lists []PublishedList `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`
	Maps  []PublishedMap  `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func NewVersion(projectId string) Version {
	return Version{
		ProjectID: projectId,
	}
}

func (Version) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.VERSION_TABLE)
}
