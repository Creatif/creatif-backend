package paginateVariables

import (
	"creatif/pkg/app/app/createProject"
	"creatif/pkg/app/domain"
	createVariable2 "creatif/pkg/app/services/variables/createVariable"
	storage2 "creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func loadEnv() {
	err := godotenv.Load("../../../../../.env")

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
	GinkgoRunSpecs(t, "Variable paginateVariables -> CRUD tests")
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
	gomega.Expect(storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLES_TABLE)).Error).Should(gomega.BeNil())
	gomega.Expect(storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.PROJECT_TABLE)).Error).Should(gomega.BeNil())
	gomega.Expect(storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.MAP_VARIABLES)).Error).Should(gomega.BeNil())
	gomega.Expect(storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLE_MAP)).Error).Should(gomega.BeNil())
})

func testCreateDeclarationVariable(projectId, name, behaviour string, groups []string, metadata []byte) createVariable2.View {
	m := map[string]interface{}{
		"one":   "one",
		"two":   []string{"one", "two", "three"},
		"three": []int{1, 2, 3},
		"four":  453,
	}

	b, err := json.Marshal(m)
	gomega.Expect(err).Should(gomega.BeNil())

	handler := createVariable2.New(createVariable2.NewModel(projectId, name, behaviour, groups, metadata, b))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	return view
}

func testCreateBasicDeclarationTextVariable(projectId, name, behaviour string) createVariable2.View {
	return testCreateDeclarationVariable(projectId, name, behaviour, []string{
		"one",
		"two",
		"three",
	}, []byte{})
}

func testAssertErrNil(err error) {
	gomega.Expect(err).Should(gomega.BeNil())
}

func testAssertIDValid(id string) {
	gomega.Expect(id).ShouldNot(gomega.BeEmpty())
	_, err := ulid.Parse(id)
	gomega.Expect(err).Should(gomega.BeNil())
}

func testCreateProject(name string) string {
	handler := createProject.New(createProject.NewModel(name))

	model, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(model.ID)

	gomega.Expect(model.Name).Should(gomega.Equal(name))

	return model.ID
}

func testAdvanceCursor(paginationId, projectId, direction, orderBy string, limit int) string {
	handler := New(NewModel(projectId, paginationId, "created_at", orderBy, direction, limit, []string{}))
	views, err := handler.Handle()
	testAssertErrNil(err)

	return views.PaginationInfo.Parameters.PaginationID
}