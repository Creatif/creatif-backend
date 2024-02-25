package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"time"
)

type Version struct {
	ID        string `gorm:"primarykey;type:text;default:gen_ulid()"`
	ProjectID string `gorm:"uniqueIndex:unique_version;type:text"`
	Name      string `gorm:"uniqueIndex:unique_version;type:text"`

	Lists      []PublishedList      `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`
	Maps       []PublishedMap       `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`
	References []PublishedReference `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func NewVersion(projectId, name string) Version {
	return Version{
		ProjectID: projectId,
		Name:      name,
	}
}

func (Version) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.VERSION_TABLE)
}
