package deleteListItemByID

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared/connections"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list item delete tests", func() {
	ginkgo.It("should delete a list item by list name and item ID", func() {
		projectId := testCreateProject("project")
		groups := testCreateGroups(projectId)
		_, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 99)
		_, referenceListId, _ := testCreateListAndReturnNameAndID(projectId, "referenceName", 100)

		var referenceListItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", referenceListId).Select("id").Find(&referenceListItems)
		testAssertErrNil(res.Error)

		addToListVariable := testAddToList(projectId, listId, []connections.Connection{
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

		handler := New(NewModel(projectId, listId, addToListVariable.ID), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(99))

		var count int
		res = storage.Gorm().Raw("SELECT count(id) AS count FROM declarations.references").Scan(&count)
		testAssertErrNil(res.Error)
		gomega.Expect(count).Should(gomega.Equal(0))
	})

	ginkgo.It("should delete a list item by list shortID and item name", func() {
		projectId := testCreateProject("project")
		testCreateGroups(projectId)
		_, listID, listShortID := testCreateListAndReturnNameAndID(projectId, "name", 100)

		var listItem declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listID).Select("short_id").First(&listItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		handler := New(NewModel(projectId, listShortID, listItem.ShortID), auth.NewTestingAuthentication(false, ""))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listID).Select("ID").Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(99))
	})
})
