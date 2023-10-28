package logger

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Logger tests", func() {
	ginkgo.It("should write info, error and warning messages to log", func() {
		logger := NewLogBuilder()
		logger.Add("one", "one")
		logger.Add("two", "two")
		logger.Add("three", "three")
		logger.Add("four", "four")
		logger.Add("five", "five")

		gomega.Expect(logger.Flush("info")).Should(gomega.BeNil())
		gomega.Expect(logger.Flush("warn")).Should(gomega.BeNil())
		gomega.Expect(logger.Flush("error")).Should(gomega.BeNil())
	})

	ginkgo.It("should write info, error and warning messages to log with same keys", func() {
		logger := NewLogBuilder()
		logger.Add("one", "one")
		logger.Add("one", "two")
		logger.Add("one", "three")
		logger.Add("one", "four")
		logger.Add("one", "five")

		gomega.Expect(logger.Flush("info")).Should(gomega.BeNil())
		gomega.Expect(logger.Flush("warn")).Should(gomega.BeNil())
		gomega.Expect(logger.Flush("error")).Should(gomega.BeNil())
	})
})
