package deleteRangeByID

import (
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list item delete tests", func() {
	ginkgo.It("should delete a range of list items by ID", func() {
		projectId := testCreateProject("project")
		listName, listId := testCreateListAndReturnNameAndID(projectId, "name", 100)

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.ListVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, listName, ids))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(90))
	})
})
