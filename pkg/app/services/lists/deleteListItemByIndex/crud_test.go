package deleteListItemByIndex

import (
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list item delete tests", func() {
	ginkgo.It("should delete a list item by index", func() {
		projectId := testCreateProject("project")
		listName, listId := testCreateListAndReturnNameAndID(projectId, "name", 100)

		var listItem declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("index").First(&listItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		handler := New(NewModel(projectId, "eng", listName, 6))
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(99))
	})
})
