package deleteRangeByID

import (
	"creatif/pkg/app/auth"
	declarations2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list item delete tests", func() {
	ginkgo.It("should delete a range of list items by name", func() {
		projectId := testCreateProject("project")
		listName, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 15)

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.ListVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, "eng", listName, "", "", ids), auth.NewNoopAuthentication(false), logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(5))
	})

	ginkgo.It("should delete a range of list items by ID", func() {
		projectId := testCreateProject("project")
		_, listId, _ := testCreateListAndReturnNameAndID(projectId, "name", 15)

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.ListVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, "eng", "", listId, "", ids), auth.NewNoopAuthentication(false), logger.NewLogBuilder())
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
		_, listId, shortID := testCreateListAndReturnNameAndID(projectId, "name", 15)

		var listItems []declarations2.ListVariable
		res := storage.Gorm().Where("list_id = ?", listId).Select("ID").Limit(10).Find(&listItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())

		ids := sdk.Map(listItems, func(idx int, value declarations2.ListVariable) string {
			return value.ID
		})

		handler := New(NewModel(projectId, "eng", "", "", shortID, ids), auth.NewNoopAuthentication(false), logger.NewLogBuilder())
		model, err := handler.Handle()
		testAssertErrNil(err)
		gomega.Expect(model).Should(gomega.BeNil())

		var remainingItems []declarations2.ListVariable
		res = storage.Gorm().Where("list_id = ?", listId).Select("ID").Find(&remainingItems)
		gomega.Expect(res.Error).Should(gomega.BeNil())
		gomega.Expect(len(remainingItems)).Should(gomega.Equal(5))
	})
})
