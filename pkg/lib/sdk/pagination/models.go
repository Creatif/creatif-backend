package pagination

import (
	"github.com/segmentio/ksuid"
	"time"
)

type InitialModel struct {
	ID        ksuid.KSUID `gorm:"primarykey;type:text CHECK(length(id)=27)"`
	CreatedAt time.Time   `gorm:"<-:create"`
}
