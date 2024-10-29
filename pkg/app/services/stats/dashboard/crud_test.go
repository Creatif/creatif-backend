package dashboard

import (
	"creatif/pkg/app/auth"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Dashboard", func() {
	ginkgo.It("Should get basic dashboard statistics", ginkgo.Label("stats", "dashboard"), func() {
		projectId := testCreateProject("project")

		testCreateMap(projectId, "map1")
		testCreateMap(projectId, "map2")
		testCreateMap(projectId, "map3")

		testCreateList(projectId, "list1")
		testCreateList(projectId, "list2")
		testCreateList(projectId, "list3")

		handler := New(NewModel(projectId), auth.NewTestingAuthentication(false, ""))
		models, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(len(models)).Should(gomega.Equal(6))

		testingRepresentation := make(map[string]string)
		testingRepresentation["list1"] = "list"
		testingRepresentation["list2"] = "list"
		testingRepresentation["list3"] = "list"
		testingRepresentation["map1"] = "map"
		testingRepresentation["map2"] = "map"
		testingRepresentation["map3"] = "map"

		for name, t := range testingRepresentation {
			found := false
			for _, model := range models {
				if model.Name == name && model.Type == t {
					found = true
					break
				}
			}
			
			gomega.Expect(found).Should(gomega.Equal(true))
		}
	})
})
