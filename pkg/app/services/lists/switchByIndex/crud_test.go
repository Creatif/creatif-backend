package switchByIndex

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"math/rand"
	"sync"
	"time"
)

var _ = ginkgo.Describe("Declaration list variable tests", func() {
	ginkgo.It("should switch two list variables indexes", func() {
		projectId := testCreateProject("project")
		indexes := testCreateListAndReturnIndexes(projectId, "list", 10)

		source := indexes[0]
		destination := indexes[5]

		handler := New(NewModel(projectId, "eng", "list", 0, 5))
		view, err := handler.Handle()
		testAssertErrNil(err)

		gomega.Expect(view.Source.Index).Should(gomega.Equal(destination))
		gomega.Expect(view.Destination.Index).Should(gomega.Equal(source))
		gomega.Expect(view.Destination.Locale).Should(gomega.Equal("eng"))
		gomega.Expect(view.Source.Locale).Should(gomega.Equal("eng"))
	})

	ginkgo.It("should switch two equal list variables indexes concurrently", func() {
		projectId := testCreateProject("project")
		testCreateListAndReturnIndexes(projectId, "list", 10)

		wg := &sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer ginkgo.GinkgoRecover()
				defer wg.Done()

				handler := New(NewModel(projectId, "eng", "list", 0, 5))
				view, err := handler.Handle()
				testAssertErrNil(err)

				gomega.Expect(view.Source.Index).ShouldNot(gomega.BeEmpty())
				gomega.Expect(view.Destination.Index).ShouldNot(gomega.BeEmpty())
				gomega.Expect(view.Destination.Locale).Should(gomega.Equal("eng"))
				gomega.Expect(view.Source.Locale).Should(gomega.Equal("eng"))
			}()
		}
		wg.Wait()
	})

	ginkgo.It("should switch two random list variables indexes concurrently", func() {
		projectId := testCreateProject("project")
		testCreateListAndReturnIndexes(projectId, "list", 10)

		randomIndex := func() int64 {
			rand.Seed(time.Now().UnixNano())
			min := 0
			max := 9
			return int64(rand.Intn(max-min+1) + min)
		}

		wg := &sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer ginkgo.GinkgoRecover()
				defer wg.Done()

				handler := New(NewModel(projectId, "eng", "list", randomIndex(), randomIndex()))
				view, err := handler.Handle()
				testAssertErrNil(err)

				gomega.Expect(view.Source.Index).ShouldNot(gomega.BeEmpty())
				gomega.Expect(view.Destination.Index).ShouldNot(gomega.BeEmpty())
				gomega.Expect(view.Destination.Locale).Should(gomega.Equal("eng"))
				gomega.Expect(view.Source.Locale).Should(gomega.Equal("eng"))
			}()
		}
		wg.Wait()
	})
})
