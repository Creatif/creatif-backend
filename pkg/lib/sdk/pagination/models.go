package pagination

import (
	"time"
)

type InitialModel struct {
	ID        string    `gorm:"primarykey;type:text CHECK(length(id)=27)"`
	CreatedAt time.Time `gorm:"<-:create"`
}
