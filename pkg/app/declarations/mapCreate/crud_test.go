package mapCreate

import (
	"creatif/pkg/app/declarations/create"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map node tests", func() {
	ginkgo.It("should create a map out of declaration nodes", func() {
		name, _ := sdk.NewULID()
		nodes := make([]create.View, 0)
		for i := 0; i < 20; i++ {
			view := testCreateBasicDeclarationTextNode(fmt.Sprintf("name-%d", i), "modifiable")

			nodes = append(nodes, view)
		}

		ids := sdk.Map(nodes, func(idx int, value create.View) string {
			return value.Name
		})

		handler := New(NewCreateMapModel(name, ids))

		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(name).Should(gomega.Equal(view.Name))
		gomega.Expect(len(view.Nodes)).Should(gomega.Equal(len(nodes)))
	})
})
