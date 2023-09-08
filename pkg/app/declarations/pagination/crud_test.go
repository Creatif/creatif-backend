package pagination

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("Declaration node pagination tests", func() {
	ginkgo.It("should return the first row of results", func() {
		for i := 0; i < 20; i++ {
			testCreateBasicAssignmentTextNode(fmt.Sprintf("name-%d", i))
		}

		handler := New(NewModel(false, "created_at", "desc", 10))
		_, err := handler.Handle()
		testAssertErrNil(err)
	})
})
