package app

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"time"
)

type Event struct {
	ID   string         `gorm:"primarykey;type:text"`
	Type string         `gorm:"index"`
	Data datatypes.JSON `gorm:"type:jsonb"`

	ProjectID string `gorm:"type:text CHECK(length(id)=27);not null"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewEvent(projectId, t string, data []byte) Event {
	return Event{
		ID:        ksuid.New().String(),
		Type:      t,
		ProjectID: projectId,
		Data:      data,
	}
}

func (Event) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.EVENTS_TABLE)
}
