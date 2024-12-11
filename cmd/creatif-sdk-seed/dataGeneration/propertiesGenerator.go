package dataGeneration

type PropertyGeneratorResult struct {
	Locale         string
	PropertyStatus string
	PropertyType   string
}

type PropertiesGenerator struct {
	Locales          []string
	PropertyStatuses []string
	PropertyTypes    []string

	currentLocaleIdx         int
	currentPropertyStatusIdx int
	currentPropertyTypeIdx   int
}

func NewPropertiesGenerator() *PropertiesGenerator {
	return &PropertiesGenerator{
		Locales:          []string{"eng", "afh", "kam", "ota", "oto"},
		PropertyStatuses: []string{"Rent", "Sell", "Rent business"},
		PropertyTypes:    []string{"House", "Apartment", "Studio", "Land"},

		currentLocaleIdx:         0,
		currentPropertyStatusIdx: 0,
		currentPropertyTypeIdx:   -1,
	}
}

/**
Index ordering:
`1. currentLocaleIdx
 2. currentPropertyStatusIdx
 3. currentPropertyTypeIdx

With every iteration, increase the currentPropertyTypeIdx. If currentPropertyTypeIdx > len(propertyTypes)-1, then increase
the currentPropertyStatusIdx by 1. If currentPropertyStatusIdx > len(propertyTypes)-1, then increase the currentLocaleIdx.

Indexes are only increased when the one below them is depleted. Upper index resets only if the one below it is depleted.
*/

func (pg *PropertiesGenerator) Generate() (PropertyGeneratorResult, bool) {
	// this is the first case only, should be executed only on first call
	if pg.currentPropertyTypeIdx == -1 && pg.currentLocaleIdx == 0 && pg.currentPropertyStatusIdx == 0 {
		pgr := PropertyGeneratorResult{
			Locale:         pg.Locales[pg.currentLocaleIdx],
			PropertyStatus: pg.PropertyStatuses[pg.currentPropertyStatusIdx],
			PropertyType:   pg.PropertyTypes[0],
		}

		pg.currentPropertyTypeIdx += 1

		return pgr, true
	}

	pg.currentPropertyTypeIdx += 1

	if pg.currentPropertyTypeIdx > len(pg.PropertyTypes)-1 {
		pg.currentPropertyStatusIdx += 1
		pg.currentPropertyTypeIdx = 0
	}

	if pg.currentPropertyStatusIdx > len(pg.PropertyStatuses)-1 {
		pg.currentLocaleIdx += 1
		pg.currentPropertyStatusIdx = 0
		pg.currentPropertyTypeIdx = 0
	}

	// when we depleted all the locales, then we stop.
	// the code that uses this generator should stop also, hence returning false
	// for every other operation below this one, this generator will return true
	if pg.currentLocaleIdx > len(pg.Locales)-1 {
		return PropertyGeneratorResult{}, false
	}

	return PropertyGeneratorResult{
		Locale:         pg.Locales[pg.currentLocaleIdx],
		PropertyStatus: pg.PropertyStatuses[pg.currentPropertyStatusIdx],
		PropertyType:   pg.PropertyTypes[pg.currentPropertyTypeIdx],
	}, true
}
