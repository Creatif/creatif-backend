package addToMap

import (
	"creatif/pkg/app/app/createProject"
	"creatif/pkg/app/domain"
	"creatif/pkg/app/services/maps/mapCreate"
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
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLES_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", domain.PROJECT_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.MAP_VARIABLES))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.VARIABLE_MAP))
})

func testCreateDetailedVariable(projectId, name, behaviour string, groups []string, metadata []byte) createVariable2.View {
	b, _ := json.Marshal(map[string]interface{}{
		"one":  1,
		"two":  "three",
		"four": "six",
	})

	handler := createVariable2.New(createVariable2.NewModel(projectId, name, behaviour, groups, metadata, b))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	return view
}

func testCreateDeclarationVariable(projectId, name, behaviour string) createVariable2.View {
	return testCreateDetailedVariable(projectId, name, behaviour, []string{}, []byte{})
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

func testCreateMap(projectId, name string, variablesNum int) mapCreate.View {
	entries := make([]mapCreate.Entry, 0)

	m := map[string]interface{}{
		"one":   "one",
		"two":   []string{"one", "two", "three"},
		"three": []int{1, 2, 3},
		"four":  453,
	}

	b, err := json.Marshal(m)
	gomega.Expect(err).Should(gomega.BeNil())

	for i := 0; i < variablesNum; i++ {
		var value interface{}
		value = "my value"
		if i%2 == 0 {
			value = true
		}

		if i%3 == 0 {
			value = map[string]interface{}{
				"one":   "one",
				"two":   []string{"one", "two", "three"},
				"three": []int{1, 2, 3},
				"four":  453,
			}
		}

		v, err := json.Marshal(value)
		gomega.Expect(err).Should(gomega.BeNil())

		variableModel := mapCreate.VariableModel{
			Name:     fmt.Sprintf("name-%d", i),
			Metadata: b,
			Groups: []string{
				"one",
				"two",
				"three",
			},
			Value:     v,
			Behaviour: "modifiable",
		}

		entries = append(entries, mapCreate.Entry{
			Type:  "variable",
			Model: variableModel,
		})
	}

	handler := mapCreate.New(mapCreate.NewModel(projectId, name, entries))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	gomega.Expect(name).Should(gomega.Equal(view.Name))
	gomega.Expect(len(view.Variables)).Should(gomega.Equal(variablesNum))

	return view
}
