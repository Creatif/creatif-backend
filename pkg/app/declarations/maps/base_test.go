package maps

import (
	"creatif/pkg/app/declarations/create"
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/appErrors"
	storage2 "creatif/pkg/lib/storage"
	"fmt"
	"github.com/google/uuid"
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
	GinkgoRunSpecs(t, "Maps -> CRUD tests")
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
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.DECLARATION_NODES_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE assignments.%s CASCADE", domain.ASSIGNMENT_NODES_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE assignments.%s CASCADE", domain.ASSIGNMENT_MAP_VALUE_NODE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE assignments.%s CASCADE", domain.ASSIGNMENT_MAP_NODES_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE assignments.%s CASCADE", domain.ASSIGNMENT_VALUE_NODE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.NODE_MAP_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.NODE_MAP_NODES_TABLE))
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

func testCreateDeclarationNode(name, behaviour string, groups []string, metadata []byte, validation create.NodeValidation) create.View {
	handler := create.New(create.NewCreateNodeModel(name, behaviour, groups, metadata, validation))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID.String())

	return view
}

func testCreateBasicDeclarationTextNode(name, behaviour string) create.View {
	return testCreateDeclarationNode(name, behaviour, []string{}, []byte{}, create.NodeValidation{})
}

func testCreateBasicDeclarationBooleanNode(name, behaviour string) create.View {
	return testCreateDeclarationNode(name, behaviour, []string{}, []byte{}, create.NodeValidation{})
}

func testAssertErrNil(err error) {
	gomega.Expect(err).Should(gomega.BeNil())
}

func testAssertIDValid(id string) {
	gomega.Expect(id).ShouldNot(gomega.BeEmpty())
	_, err := uuid.Parse(id)
	gomega.Expect(err).Should(gomega.BeNil())
}
