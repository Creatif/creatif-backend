package deleteListItemByID

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list item delete tests", func() {
	ginkgo.It("should delete a list item by list name and item ID", func() {
		projectId := testCreateProject("project")
		listName, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 100)

		var listItem declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("ID").First(&listItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		handler := New(NewModel(projectId, "eng", listName, "", "", listItem.ID, ""), auth.NewNoopAuthentication(false), logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(99))
	})

	ginkgo.It("should delete a list item by list shortID and item name", func() {
		projectId := testCreateProject("project")
		_, listID, listShortID := testCreateListAndReturnNameAndID(projectId, "name", 100)

		var listItem declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listID).Select("short_id").First(&listItem)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		handler := New(NewModel(projectId, "eng", "", "", listShortID, "", listItem.ShortID), auth.NewNoopAuthentication(false), logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var listItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listID).Select("ID").Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(listItems)).Should(gomega.Equal(99))
	})
})
