package app

import (
	"creatif/pkg/app/domain"
	"fmt"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
	"time"
)

type Activity struct {
	ID   string         `gorm:"primarykey;type:text"`
	Data datatypes.JSON `gorm:"type:jsonb"`

	CreatedAt time.Time `gorm:"<-:create"`
}

func NewActivity(data datatypes.JSON) Activity {
	return Activity{
		ID:   ksuid.New().String(),
		Data: data,
	}
}

func (Activity) TableName() string {
	return fmt.Sprintf("%s.%s", "app", domain.ACTIVITY)
}
