package declarations

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"time"
)

type File struct {
	ID        string    `gorm:"primarykey;type:text"`
	ListID    *string   `gorm:"type:text;default null"`
	MapID     *string   `gorm:"type:text;default null"`
	ProjectID string    `gorm:"primarykey;type:text"`
	Name      string    `gorm:"type:text"`
	FileName  string    `gorm:"type:text"`
	FieldName string    `gorm:"type:text"`
	MimeType  string    `gorm:"type:text"`
	Extension string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewFile(projectId string, listId *string, mapId *string, name, fieldName, mimeType, extension, fileName string) File {
	return File{
		ID:        ksuid.New().String(),
		FieldName: fieldName,
		FileName:  fileName,
		ListID:    listId,
		MapID:     mapId,
		Name:      name,
		Extension: extension,
		MimeType:  mimeType,
		ProjectID: projectId,
	}
}

func (File) TableName() string {
	return fmt.Sprintf("%s.%s", "declarations", domain.FILE_TABLE)
}
