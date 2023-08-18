package create

import (
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/appErrors"
	storage2 "creatif/pkg/lib/storage"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func loadEnv() {
	err := godotenv.Load("../../../../.env")

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
	GinkgoRunSpecs(t, "Declaration -> CRUD tests")
}

var _ = ginkgo.BeforeSuite(func() {
	loadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Zagreb",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	err := storage2.Connect(dsn)

	if err != nil {
		log.Fatalln(err)
	}
})

var _ = GinkgoAfterSuite(func() {
	sql, err := storage2.SQLDB()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("Could not retreive DB instance: %s", err))
	}

	if err := sql.Close(); err != nil {
		ginkgo.Fail(fmt.Sprintf("Could not close database connection: %s", err))
	}
})

var _ = GinkgoAfterHandler(func() {
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", domain.DECLARATION_NODES_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", domain.ASSIGNMENT_NODES_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", domain.ASSIGNMENT_NODE_BOOLEAN_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", domain.ASSIGNMENT_NODE_TEXT_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", domain.USERS_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", domain.PROJECT_TABLE))
})

func _assertValidation(err error, keys []string) {
	validationError, ok := err.(appErrors.AppError[map[string]string])
	if ok {
		data := validationError.Data()

		for key := range data {
			gomega.Expect(keys).Should(gomega.ContainElement(key))
		}
	}
}