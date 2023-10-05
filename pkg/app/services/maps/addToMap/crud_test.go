package addToMap

import (
	getMap2 "creatif/pkg/app/services/maps/getMap"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (UPDATE) map entry tests", func() {
	ginkgo.It("should add an entry to the map", func() {
		projectId := testCreateProject("project")
		m := testCreateMap(projectId, "mapName", 10)

		handler := New(NewModel(projectId, m.Name, VariableModel{
			Name:      "newEntry",
			Metadata:  nil,
			Groups:    nil,
			Behaviour: "readonly",
			Value:     nil,
		}))

		_, err := handler.Handle()
		testAssertErrNil(err)

		getMapHandler := getMap2.New(getMap2.NewModel(projectId, m.Name, []string{}))
		maps, err := getMapHandler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(maps.ID)
		gomega.Expect(maps.ProjectID).Should(gomega.Equal(projectId))

		gomega.Expect(len(maps.Variables)).Should(gomega.Equal(11))

		found := false
		for _, variable := range maps.Variables {
			if variable["name"].(string) == "newEntry" {
				found = true
			}
		}

		gomega.Expect(found).Should(gomega.BeTrue())
	})
})