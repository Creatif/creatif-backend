package published

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"time"
)

type PublishedFile struct {
	ID        string    `gorm:"primarykey;type:text"`
	VersionID string    `gorm:"primarykey:type:text"`
	ProjectID string    `gorm:"primarykey;type:text"`
	Name      string    `gorm:"type:text"`
	FileName  string    `gorm:"type:text"`
	MimeType  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewFile(projectId string, name, mimeType, fileName string) PublishedFile {
	return PublishedFile{
		ID:        ksuid.New().String(),
		FileName:  fileName,
		Name:      name,
		MimeType:  mimeType,
		ProjectID: projectId,
	}
}

func (PublishedFile) TableName() string {
	return fmt.Sprintf("%s.%s", "published", domain.PUBLISHED_FILES_TABLE)
}
