package deleteListItemByID

import (
	"creatif/pkg/app/app/createProject"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain"
	createList2 "creatif/pkg/app/services/lists/createList"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/logger"
	storage2 "creatif/pkg/lib/storage"
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
	GinkgoRunSpecs(t, "Declaration Lists -> CRUD tests")
}

func runLogger() {
	if err := logger.BuildLoggers(os.Getenv("LOG_DIRECTORY")); err != nil {
		log.Fatalln(fmt.Sprintf("Cannot createProject logger: %s", err.Error()))
	}

	logger.Info("Health info logger health check... Ignore!")
	logger.Warn("Health warning logger health check... Ignore!")
	logger.Error("Health error logger health check... Ignore!")
}

var _ = ginkgo.BeforeSuite(func() {
	loadEnv()
	runLogger()

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

	gomega.Expect(locales.Store()).Should(gomega.BeNil())
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
	res := storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLES_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.PROJECT_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.MAP_VARIABLES))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLE_MAP))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.LIST_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.LIST_VARIABLES_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
	res = storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.USERS_TABLE))
	gomega.Expect(res.Error).Should(gomega.BeNil())
})

func testAssertErrNil(err error) {
	gomega.Expect(err).Should(gomega.BeNil())
}

func testAssertIDValid(id string) {
	gomega.Expect(id).ShouldNot(gomega.BeEmpty())
	_, err := ulid.Parse(id)
	gomega.Expect(err).Should(gomega.BeNil())
}

func testCreateProject(name string) string {
	handler := createProject.New(createProject.NewModel(name), auth.NewNoopAuthentication(), logger.NewLogBuilder())

	model, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(model.ID)

	gomega.Expect(model.Name).Should(gomega.Equal(name))

	return model.ID
}

func testCreateListAndReturnNameAndID(projectId, name string, varNum int) (string, string) {
	variables := make([]createList2.Variable, varNum)
	for i := 0; i < varNum; i++ {
		variables[i] = createList2.Variable{
			Name:      fmt.Sprintf("one-%d", i),
			Metadata:  nil,
			Groups:    nil,
			Behaviour: "readonly",
			Value:     nil,
		}
	}

	handler := createList2.New(createList2.NewModel(projectId, "eng", name, variables), auth.NewNoopAuthentication(), logger.NewLogBuilder())

	list, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(list.ID)

	gomega.Expect(list.Name).Should(gomega.Equal(name))

	return list.Name, list.ID
}
