package switchByID

import (
	"creatif/pkg/lib/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"math/rand"
	"sync"
	"time"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should switch two list variables indexes", func() {
		projectId := testCreateProject("project")
		idsAndIndexes := testCreateListAndReturnIdsAndIndexes(projectId, "list", 10)

		source := idsAndIndexes[0]
		destination := idsAndIndexes[5]

		handler := New(NewModel(projectId, "eng", "list", source["id"], destination["id"]), logger.NewLogBuilder())
		view, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(view.Destination.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(view.Source.Locale).Should(gomega.Equal("eng"))
	})

	ginkgo.It("should switch two equal list variables indexes concurrently", func() {
		projectId := testCreateProject("project")
		ids := testCreateListAndReturnIdsAndIndexes(projectId, "list", 10)

		source := ids[0]
		destination := ids[5]

		wg := &sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer ginkgo.GinkgoRecover()
				defer wg.Done()

				handler := New(NewModel(projectId, "eng", "list", source["id"], destination["id"]), logger.NewLogBuilder())
				view, err := handler.Handle()
				testAssertErrNil(err)

				gomega.Expect(view.Destination.Locale).Should(gomega.Equal("eng"))
				gomega.Expect(view.Source.Locale).Should(gomega.Equal("eng"))
			}()
		}
		wg.Wait()
	})

	ginkgo.It("should switch two random list variables indexes concurrently", func() {
		projectId := testCreateProject("project")
		ids := testCreateListAndReturnIdsAndIndexes(projectId, "list", 10)

		randomIndex := func() int {
			rand.Seed(time.Now().UnixNano())
			min := 0
			max := 9
			return rand.Intn(max-min+1) + min
		}

		sourceIdx := 0
		destinationIdx := 0
		for {
			one := randomIndex()
			two := randomIndex()
			if one != two {
				sourceIdx = one
				destinationIdx = two

				break
			}
		}

		wg := &sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer ginkgo.GinkgoRecover()
				defer wg.Done()

				handler := New(NewModel(projectId, "eng", "list", ids[sourceIdx]["id"], ids[destinationIdx]["id"]), logger.NewLogBuilder())
				view, err := handler.Handle()
				testAssertErrNil(err)

				gomega.Expect(view.Destination.Locale).Should(gomega.Equal("eng"))
				gomega.Expect(view.Source.Locale).Should(gomega.Equal("eng"))
			}()
		}
		wg.Wait()
	})
})
