package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

type Version struct {
	ID                  string `gorm:"primarykey;type:text"`
	ProjectID           string `gorm:"uniqueIndex:unique_version;type:text"`
	Name                string `gorm:"uniqueIndex:unique_version;type:text"`
	IsProductionVersion bool

	Lists      []PublishedList      `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`
	Maps       []PublishedMap       `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`
	References []PublishedReference `gorm:"foreignKey:VersionID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func (u *Version) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = ksuid.New().String()
	
	return nil
}

func NewVersion(projectId, name string, isProduction bool) Version {
	return Version{
		ProjectID:           projectId,
		Name:                name,
		IsProductionVersion: isProduction,
	}
}

func (Version) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.VERSION_TABLE)
}
