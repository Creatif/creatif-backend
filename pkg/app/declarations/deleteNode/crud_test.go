package deleteNode

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("Declaration (DELETE) node tests", func() {
	ginkgo.It("should delete a declaration node and all assignment nodes", func() {
		view := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewModel(view.Name))

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("id = ?", view.ID).First(&declarations.Node{})
		gomega.Expect(errors.Is(res.Error, gorm.ErrRecordNotFound)).Should(gomega.BeTrue())
	})
})
