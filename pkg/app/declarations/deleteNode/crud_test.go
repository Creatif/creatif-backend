package deleteNode

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/lib/storage"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration (DELETE) node tests", func() {
	ginkgo.It("should delete a declaration node and all assignment nodes", func() {
		view := testCreateBasicDeclarationTextNode("node", "modifiable")

		handler := New(NewModel(view.Name))

		_, err := handler.Handle()
		testAssertErrNil(err)

		res := storage.Gorm().Where("declaration_node_id = ?", view.ID).First(&assignments.Node{})
		gomega.Expect(res.Error).ShouldNot(gomega.BeNil())
	})
})
