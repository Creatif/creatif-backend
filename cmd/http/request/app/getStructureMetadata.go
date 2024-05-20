package app

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type GetStructureMetadataConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type GetStructureMetadata struct {
	ID     string                       `param:"projectId"`
	Config []GetStructureMetadataConfig `json:"config"`
}

func SanitizeGetStructureMetadata(model GetStructureMetadata) GetStructureMetadata {
	p := bluemonday.StrictPolicy()
	model.ID = p.Sanitize(model.ID)
	model.Config = sdk.Map(model.Config, func(idx int, value GetStructureMetadataConfig) GetStructureMetadataConfig {
		return GetStructureMetadataConfig{
			Name: p.Sanitize(value.Name),
			Type: p.Sanitize(value.Type),
		}
	})

	return model
}
