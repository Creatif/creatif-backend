package dashboard

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publishing/publish"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Dashboard", func() {
	ginkgo.It("Should get basic dashboard statistics", ginkgo.Label("stats", "dashboard"), func() {
		projectId := testCreateProject("project")

		view1 := testCreateMap(projectId, "map1")
		for i := 0; i < 100; i++ {
			testAddToMap(projectId, view1.ID, fmt.Sprintf("name-%d", i), nil, nil)
		}

		view2 := testCreateMap(projectId, "map2")
		for i := 0; i < 100; i++ {
			testAddToMap(projectId, view2.ID, fmt.Sprintf("name-%d", i), nil, nil)
		}

		view3 := testCreateMap(projectId, "map3")
		for i := 0; i < 100; i++ {
			testAddToMap(projectId, view3.ID, fmt.Sprintf("name-%d", i), nil, nil)
		}

		view4 := testCreateList(projectId, "list1")
		for i := 0; i < 100; i++ {
			testAddToList(projectId, view4.ID, fmt.Sprintf("name-%d", i), nil, nil)
		}

		view5 := testCreateList(projectId, "list2")
		for i := 0; i < 100; i++ {
			testAddToList(projectId, view5.ID, fmt.Sprintf("name-%d", i), nil, nil)
		}

		view6 := testCreateList(projectId, "list3")
		for i := 0; i < 100; i++ {
			testAddToList(projectId, view6.ID, fmt.Sprintf("name-%d", i), nil, nil)
		}

		// publish multiple versions
		for i := 0; i < 5; i++ {
			handler := publish.New(publish.NewModel(projectId, fmt.Sprintf("v%d", i)), auth.NewTestingAuthentication(false, ""))
			model, err := handler.Handle()
			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(model.ID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(model.Name).ShouldNot(gomega.BeEmpty())
		}

		handler := New(NewModel(projectId), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(len(model.Structures)).Should(gomega.Equal(6))
		gomega.Expect(len(model.Versions)).Should(gomega.Equal(5))

		for _, structure := range model.Structures {
			gomega.Expect(structure.Name).ShouldNot(gomega.BeEmpty())
			gomega.Expect(structure.Count).Should(gomega.Equal(100))
			gomega.Expect(structure.ID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(structure.Type).ShouldNot(gomega.BeEmpty())
			gomega.Expect(structure.CreatedAt).ShouldNot(gomega.BeEmpty())
		}

		testingRepresentation := make(map[string]string)
		testingRepresentation["list1"] = "list"
		testingRepresentation["list2"] = "list"
		testingRepresentation["list3"] = "list"
		testingRepresentation["map1"] = "map"
		testingRepresentation["map2"] = "map"
		testingRepresentation["map3"] = "map"

		for name, t := range testingRepresentation {
			found := false
			for _, model := range model.Structures {
				if model.Name == name && model.Type == t {
					found = true
					break
				}
			}

			gomega.Expect(found).Should(gomega.Equal(true))
		}

		versions := model.Versions
		for _, version := range versions {
			gomega.Expect(version.Name).ShouldNot(gomega.BeEmpty())
			gomega.Expect(version.ID).ShouldNot(gomega.BeEmpty())
			gomega.Expect(version.CreatedAt).ShouldNot(gomega.BeEmpty())
		}
	})
})
