package publish

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/groups/addGroups"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"os"
)

var _ = ginkgo.Describe("Publishing", func() {
	ginkgo.It("Should publish all lists and maps", ginkgo.Label("publish", "create"), func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId, 5)

		map1 := testCreateMap(projectId, "map1")
		map2 := testCreateMap(projectId, "map2")
		map3 := testCreateMap(projectId, "map3")

		list1 := testCreateList(projectId, "list1")
		list2 := testCreateList(projectId, "list2")
		list3 := testCreateList(projectId, "list3")

		referenceMap := testCreateMap(projectId, "referenceMap")
		referenceMapItem1 := testAddToMap(projectId, referenceMap.ID, "reference-map-1", []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))
		referenceMapItem2 := testAddToMap(projectId, referenceMap.ID, "reference-map-2", []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		referenceList := testCreateList(projectId, "referenceList")
		referenceListItem1 := testAddToList(projectId, referenceList.ID, "reference-list-1", []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))
		referenceListItem2 := testAddToList(projectId, referenceList.ID, "reference-list-2", []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
			return value.ID
		}))

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, map1.ID, fmt.Sprintf("map-%d", i), []connections.Connection{
				{
					Path:          "first",
					StructureType: "map",
					VariableID:    referenceMapItem1.Variable.ID,
				},
				{
					Path:          "second",
					StructureType: "map",
					VariableID:    referenceMapItem2.Variable.ID,
				},
			}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, map2.ID, fmt.Sprintf("map-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		for i := 0; i < 100; i++ {
			testAddToMap(projectId, map3.ID, fmt.Sprintf("map-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		for i := 0; i < 100; i++ {
			testAddToList(projectId, list1.ID, fmt.Sprintf("list-%d", i), []connections.Connection{
				{
					Path:          "first",
					StructureType: "list",
					VariableID:    referenceListItem1.ID,
				},
				{
					Path:          "second",
					StructureType: "list",
					VariableID:    referenceListItem2.ID,
				},
			}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		for i := 100; i < 200; i++ {
			testAddToList(projectId, list2.ID, fmt.Sprintf("list-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		for i := 200; i < 300; i++ {
			testAddToList(projectId, list3.ID, fmt.Sprintf("list-%d", i), []connections.Connection{}, sdk.Map(groups, func(idx int, value addGroups.View) string {
				return value.ID
			}))
		}

		handler := New(NewModel(projectId, "version name"), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(model.ID).ShouldNot(gomega.BeEmpty())
		gomega.Expect(model.Name).Should(gomega.Equal("version name"))

		var versionCount int64
		res := storage.Gorm().Raw("SELECT count(*) FROM published.versions").Scan(&versionCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(versionCount).Should(gomega.Equal(int64(1)))

		var listsCount int64
		res = storage.Gorm().Raw("SELECT count(*) FROM published.published_lists").Scan(&listsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(listsCount).Should(gomega.Equal(int64(302)))

		var mapsCount int64
		res = storage.Gorm().Raw("SELECT count(*) FROM published.published_maps").Scan(&mapsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(mapsCount).Should(gomega.Equal(int64(302)))

		var groupsCount int64
		res = storage.Gorm().Raw("SELECT count(*) FROM published.published_groups").Scan(&groupsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(mapsCount).Should(gomega.Equal(int64(302)))

		var connectionsCount int64
		res = storage.Gorm().Raw("SELECT count(*) FROM published.published_connections").Scan(&connectionsCount)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(connectionsCount).Should(gomega.Equal(int64(400)))

		fileInfo, err := os.Stat(fmt.Sprintf("/app/public/%s/%s", projectId, "version name"))
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(fileInfo.IsDir()).Should(gomega.BeTrue())
	})
})
