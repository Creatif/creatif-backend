package deleteVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("Declaration (DELETE) variable tests", func() {
	ginkgo.It("should delete a declaration variable and all assignment variables", func() {
		projectId := testCreateProject("project")
		view := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, view.Name))

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("id = ?", view.ID).First(&declarations.Variable{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})
})
