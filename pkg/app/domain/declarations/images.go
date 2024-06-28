package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"time"
)

type Image struct {
	ID          string    `gorm:"primarykey;type:text"`
	StructureID string    `gorm:"primarykey;type:text"`
	ProjectID   string    `gorm:"primarykey;type:text"`
	Name        string    `gorm:"type:text"`
	FieldName   string    `gorm:"type:text"`
	MimeType    string    `gorm:"type:text"`
	Extension   string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
}

func NewImage(projectId, structureId, name, fieldName, mimeType, extension string) Image {
	return Image{
		ID:          ksuid.New().String(),
		FieldName:   fieldName,
		StructureID: structureId,
		Name:        name,
		Extension:   extension,
		MimeType:    mimeType,
		ProjectID:   projectId,
	}
}

func (Image) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.IMAGE_TABLE)
}
