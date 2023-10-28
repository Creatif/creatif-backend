package logger

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func loadEnv() {
	err := godotenv.Load("../../../.env")

	if err != nil {
		log.Fatal(err)
	}
}

var GomegaRegisterFailHandler = gomega.RegisterFailHandler
var GinkgoFail = ginkgo.Fail
var GinkgoRunSpecs = ginkgo.RunSpecs
var GinkgoAfterHandler = ginkgo.AfterEach
var GinkgoAfterSuite = ginkgo.AfterSuite

func TestApi(t *testing.T) {
	GomegaRegisterFailHandler(GinkgoFail)
	GinkgoRunSpecs(t, "Logger tests")
}

var _ = ginkgo.BeforeSuite(func() {
	loadEnv()
	err := os.RemoveAll(os.Getenv("LOG_DIRECTORY"))
	if err != nil {
		log.Fatal(err)
	}

	runLogger()
})

var _ = GinkgoAfterSuite(func() {
	err := os.RemoveAll(os.Getenv("LOG_DIRECTORY"))
	if err != nil {
		log.Fatal(err)
	}
})

func runLogger() {
	if err := BuildLoggers(os.Getenv("LOG_DIRECTORY")); err != nil {
		log.Fatalln(fmt.Sprintf("Cannot createProject logger: %s", err.Error()))
	}

	Info("Health info logger health check... Ignore!")
	Warn("Health warning logger health check... Ignore!")
	Error("Health error logger health check... Ignore!")
}
