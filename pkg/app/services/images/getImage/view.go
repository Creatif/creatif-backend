package getImage

import (
	"creatif/pkg/app/domain/declarations"
)

type View struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MimeType  string `json:"mimeType"`
	FilePath  string `json:"filePath"`
	Extension string `json:"extension"`
}

func newView(model declarations.Image) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		MimeType:  model.MimeType,
		FilePath:  model.Name,
		Extension: model.Extension,
	}
}
