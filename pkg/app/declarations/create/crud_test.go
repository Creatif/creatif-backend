package create

import (
	"creatif/pkg/app/domain/assignments"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Declaration node tests", func() {
	ginkgo.It("should create a text declaration node", func() {
		name := uuid.NewString()
		handler := New(NewCreateNodeModel(name, assignments.ValueTextType, "modifiable", []string{}, []byte{}, NodeValidation{}))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Type).Should(gomega.Equal(assignments.ValueTextType))
	})

	ginkgo.It("should create a boolean declaration node", func() {
		name := uuid.NewString()
		handler := New(NewCreateNodeModel(name, assignments.ValueBooleanType, "modifiable", []string{}, []byte{}, NodeValidation{}))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		gomega.Expect(view.Name).ShouldNot(gomega.BeEmpty())
		gomega.Expect(view.Behaviour).Should(gomega.Equal("modifiable"))
		gomega.Expect(view.Type).Should(gomega.Equal(assignments.ValueBooleanType))
	})
})
