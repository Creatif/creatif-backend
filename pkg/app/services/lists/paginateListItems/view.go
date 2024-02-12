package paginateListItems

import (
	"creatif/pkg/app/services/locales"
	"time"
)

type View struct {
	ID        string      `json:"id"`
	ShortID   string      `json:"shortId"`
	Locale    string      `json:"locale"`
	Index     float64     `json:"index"`
	Name      string      `json:"name"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []QueryVariable) ([]View, error) {
	views := make([]View, len(models))
	for i, value := range models {
		locale, err := locales.GetAlphaWithID(value.LocaleID)
		if err != nil {
			return nil, err
		}

		views[i] = View{
			ID:        value.ID,
			Name:      value.Name,
			Index:     value.Index,
			Locale:    locale,
			ShortID:   value.ShortID,
			Groups:    value.Groups,
			Value:     value.Value,
			Behaviour: value.Behaviour,
			Metadata:  value.Metadata,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	}

	return views, nil
}
