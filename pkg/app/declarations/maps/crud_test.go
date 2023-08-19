package create

import (
	"creatif/pkg/app/declarations/create"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Map node tests", func() {
	ginkgo.It("should create a map out of declaration nodes", func() {
		nodes := make([]create.View, 0)
		for i := 0; i < 20; i++ {
			view := testCreateBasicDeclarationTextNode(fmt.Sprintf("name-%d", i), "this is a text node")

			nodes = append(nodes, view)
		}

		ids := sdk.Map(nodes, func(idx int, value create.View) string {
			return value.ID
		})

		handler := New(NewCreateMapModel(ids))

		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())

		gomega.Expect(len(view.Names)).Should(gomega.Equal(len(nodes)))
	})
})
