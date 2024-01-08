package appendToList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"github.com/lib/pq"
)

func assignDefaultGroupsToVariables(variables []Variable) {
	for _, v := range variables {
		pqGroups := pq.StringArray{}
		if v.Groups == nil {
			v.Groups = pq.StringArray{}
		}

		for _, k := range v.Groups {
			pqGroups = append(pqGroups, k)
		}

		v.Groups = pqGroups
	}
}

func createListVariables(listID string, variables []Variable) ([]declarations.ListVariable, error) {
	highestIndex, err := getHighestIndex(listID)
	if err != nil {
		return []declarations.ListVariable{}, err
	}

	listVariables := make([]declarations.ListVariable, len(variables))
	for i := 0; i < len(variables); i++ {
		if variables[i].Locale == "" {
			variables[i].Locale = "eng"
		}

		localeID, _ := locales.GetIDWithAlpha(variables[i].Locale)
		v := variables[i]
		listVariables[i] = declarations.NewListVariable(listID, localeID, v.Name, v.Behaviour, v.Metadata, v.Groups, v.Value)
		listVariables[i].Index = float64(highestIndex + 1000)
		highestIndex += 1000
	}

	return listVariables, nil
}
