package createAndDiff

import "creatif/pkg/lib/sdk"

func processMetadata(metadata []MetadataModel) PreViewModel {
	preViewModel := PreViewModel{
		ID:    metadata[0].ID,
		Name:  metadata[0].Name,
		State: metadata[0].State,
		Maps:  make([]PreViewStructure, 0),
		Lists: make([]PreViewStructure, 0),
	}

	if len(metadata) == 1 && metadata[0].Map == "" && metadata[0].List == "" {
		return preViewModel
	}

	for _, v := range metadata {
		if v.Map != "" {
			if !sdk.IncludesFn(preViewModel.Maps, func(item PreViewStructure) bool {
				return item.Name == v.Map
			}) {
				preViewModel.Maps = append(preViewModel.Maps, PreViewStructure{
					Name:    v.Map,
					ID:      v.MapID,
					ShortID: v.MapShortID,
				})
			}
		}

		if v.List != "" {
			if !sdk.IncludesFn(preViewModel.Lists, func(item PreViewStructure) bool {
				return item.Name == v.List
			}) {
				preViewModel.Lists = append(preViewModel.Lists, PreViewStructure{
					Name:    v.List,
					ID:      v.ListID,
					ShortID: v.ListShortID,
				})
			}
		}
	}

	return preViewModel
}
