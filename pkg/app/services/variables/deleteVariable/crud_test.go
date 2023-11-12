package deleteVariable

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("Declaration (DELETE) variable tests", func() {
	ginkgo.It("should delete a variable by name", func() {
		projectId := testCreateProject("project")
		view := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, "", "", view.Name, "eng"), auth.NewNoopAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("id = ?", view.ID).First(&declarations.Variable{})
		gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})

	ginkgo.It("should delete a variable by ID", func() {
		projectId := testCreateProject("project")
		view := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, view.ID, "", "", "eng"), auth.NewNoopAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("id = ?", view.ID).First(&declarations.Variable{})
		gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})

	ginkgo.It("should delete a variable by shortID", func() {
		projectId := testCreateProject("project")
		view := testCreateDeclarationVariable(projectId, "variable", "modifiable")

		handler := New(NewModel(projectId, "", view.ShortID, "", "eng"), auth.NewNoopAuthentication(false), logger.NewLogBuilder())

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("id = ?", view.ID).First(&declarations.Variable{})
		gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})
})
