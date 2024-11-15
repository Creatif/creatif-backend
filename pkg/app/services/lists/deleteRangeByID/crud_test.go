package deleteRangeByID

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/lists/addToList"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list item delete tests", func() {
	ginkgo.It("should delete a range of list items by name", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId)
		listName, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 15)
		_, referenceListId, _ := testCreateListAndReturnNameAndID(projectId, "referenceList", 10)

		var referenceListItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", referenceListId).Select("id").Find(&referenceListItems)
		testAssertErrNil(res.Error)

		addedMapsWithReferences := make([]addToList.View, 0)
		for i := 0; i < 10; i++ {
			addToMapVariable := testAddToList(projectId, listId, fmt.Sprintf("newAdd-%d", i), []connections.Connection{
				{
					Path:          "first",
					StructureType: "list",
					VariableID:    referenceListItems[0].ID,
				},
				{
					Path:          "second",
					StructureType: "list",
					VariableID:    referenceListItems[1].ID,
				},
			}, groups)

			addedMapsWithReferences = append(addedMapsWithReferences, addToMapVariable)
		}

		var listItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.ListVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, listName, ids), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(0))

		var count int
		res = storage.Gorm().Raw("SELECT count(child_variable_id) AS count FROM declarations.connections").Scan(&count)
		testAssertErrNil(res.Error)
		gomega.Expect(count).Should(gomega.Equal(0))
	})

	ginkgo.It("should delete a range of list items by ID", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId)
		_, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 15)

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.ListVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, listId, ids), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(5))
	})

	ginkgo.It("should delete a range of list items by shortID", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId)
		_, listId, shortID := testCreateListAndReturnNameAndID(projectId, "name", 15)

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.ListVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, shortID, ids), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(5))
	})
})
