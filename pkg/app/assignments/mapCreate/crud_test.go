package mapCreate

import (
	"creatif/pkg/app/declarations/create"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Assignment map CRUD tests", func() {
	ginkgo.It("should assign a value to a map", func() {
		nodes := make([]create.View, 0)
		for i := 0; i < 10; i++ {
			nodes = append(nodes, testCreateBasicDeclarationTextNode(fmt.Sprintf("name-%d", i), "modifiable"))
		}

		m := testCreateMap("mapName", sdk.Map(nodes, func(idx int, value create.View) string {
			return value.ID
		}))

		b, _ := json.Marshal("this is a text value")
		handler := New(NewAssignValueModel(m.Name, b))

		view, err := handler.Handle()
		gomega.Expect(err).Should(gomega.BeNil())
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var value string
		testAssertErrNil(json.Unmarshal(view.Value.([]byte), &value))
		gomega.Expect(value).Should(gomega.Equal("this is a text value"))
	})
})
