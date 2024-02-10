package getListGroups

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should get all distinct groups from a list", func() {
		projectId := testCreateProject("project")
		view := testCreateList(projectId, "list", 5)

		var itemId string
		res := storage.Gorm().Raw("SELECT id FROM declarations.list_variables WHERE list_id = ? LIMIT 1", view.ID).Scan(&itemId)
		testAssertErrNil(res.Error)

		l := logger.NewLogBuilder()
		handler := New(NewModel(view.ID, itemId, projectId), auth.NewTestingAuthentication(true, ""), l)
		groups, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(len(groups)).To(gomega.Equal(3))
	})
})
