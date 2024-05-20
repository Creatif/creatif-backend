package createAndDiff

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
)

type View struct {
	Metadata   ViewModel       `json:"metadata"`
	Diff       ViewDiff        `json:"diff"`
	Structures []ViewStructure `json:"structures"`
}

type ViewStructure struct {
	Name          string `json:"name"`
	ID            string `json:"id"`
	ShortID       string `json:"shortId"`
	StructureType string `json:"structureType"`
}

type ViewList struct {
	ID      string `json:"id"`
	ShortID string `json:"shortId"`
	Name    string `json:"name"`
}

type ViewMap struct {
	ID      string `json:"id"`
	ShortID string `json:"shortId"`
	Name    string `json:"name"`
}

type ViewModel struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	State string          `json:"state"`
	Maps  []ViewStructure `json:"maps"`
	Lists []ViewStructure `json:"lists"`
}

type ViewDiff struct {
	Lists []ViewList `json:"lists"`
	Maps  []ViewMap  `json:"maps"`
}

func newView(model LogicModel) View {
	viewModel := ViewModel{
		ID:    model.Metadata.ID,
		Name:  model.Metadata.Name,
		State: model.Metadata.State,
		Maps: sdk.Map(model.Metadata.Maps, func(idx int, value PreViewStructure) ViewStructure {
			return ViewStructure{
				ID:      value.ID,
				ShortID: value.ShortID,
				Name:    value.Name,
			}
		}),
		Lists: sdk.Map(model.Metadata.Lists, func(idx int, value PreViewStructure) ViewStructure {
			return ViewStructure{
				ID:      value.ID,
				ShortID: value.ShortID,
				Name:    value.Name,
			}
		}),
	}

	diff := ViewDiff{
		Lists: sdk.Map(model.Diff.Lists, func(idx int, value declarations.List) ViewList {
			return ViewList{
				ID:      value.ID,
				ShortID: value.ShortID,
				Name:    value.Name,
			}
		}),
		Maps: sdk.Map(model.Diff.Maps, func(idx int, value declarations.Map) ViewMap {
			return ViewMap{
				ID:      value.ID,
				ShortID: value.ShortID,
				Name:    value.Name,
			}
		}),
	}

	structures := sdk.Map(model.Structures, func(idx int, value ListOrMap) ViewStructure {
		return ViewStructure{
			Name:          value.Name,
			ID:            value.ID,
			ShortID:       value.ShortID,
			StructureType: value.StructureType,
		}
	})

	return View{
		Metadata:   viewModel,
		Diff:       diff,
		Structures: structures,
	}
}
