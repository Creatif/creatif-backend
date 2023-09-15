package getBatchNodes

import (
	assignmentsCreate "creatif/pkg/app/assignments/create"
	"creatif/pkg/app/declarations/createNode"
	mapsCreate "creatif/pkg/app/declarations/mapCreate"
	"creatif/pkg/app/domain"
	"creatif/pkg/lib/sdk"
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
	GinkgoRunSpecs(t, "Combined -> getBatchNodes tests")
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
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE assignments.%s CASCADE", domain.ASSIGNMENT_VALUE_NODE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.NODE_MAP_TABLE))
	storage2.Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", domain.NODE_MAP_NODES_TABLE))
})

func testUniqueName() string {
	uid, err := sdk.NewULID()
	gomega.Expect(err).Should(gomega.BeNil())
	return uid
}

func testCreateDeclarationNode(name, behaviour string, groups []string, metadata []byte) createNode.View {
	handler := createNode.New(createNode.NewModel(name, behaviour, groups, metadata))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	return view
}

func testCreateBasicDeclarationTextNode(name, behaviour string) createNode.View {
	return testCreateDeclarationNode(name, behaviour, []string{
		"one",
		"two",
		"three",
	}, []byte{})
}

func testCreateBasicDeclarationBooleanNode(name, behaviour string) createNode.View {
	return testCreateDeclarationNode(name, behaviour, []string{}, []byte{})
}

func testCreateMap(name string, nodesNum int) mapsCreate.View {
	entries := make([]mapsCreate.Entry, 0)

	m := map[string]interface{}{
		"one":   "one",
		"two":   []string{"one", "two", "three"},
		"three": []int{1, 2, 3},
		"four":  453,
	}

	b, err := json.Marshal(m)
	gomega.Expect(err).Should(gomega.BeNil())

	for i := 0; i < nodesNum; i++ {
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

		nodeModel := mapsCreate.NodeModel{
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

		entries = append(entries, mapsCreate.Entry{
			Type:  "node",
			Model: nodeModel,
		})
	}

	handler := mapsCreate.New(mapsCreate.NewModel(name, entries))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	gomega.Expect(name).Should(gomega.Equal(view.Name))
	gomega.Expect(len(view.Nodes)).Should(gomega.Equal(nodesNum))

	return view
}

func testCreateBasicAssignmentTextNode(name string, value interface{}) assignmentsCreate.View {
	declarationNode := testCreateBasicDeclarationTextNode(name, "modifiable")

	b, _ := json.Marshal(value)

	handler := assignmentsCreate.New(assignmentsCreate.NewCreateNodeModel(declarationNode.Name, b))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	return view
}

func testCreateBasicAssignmentBooleanNode(name string, value bool) assignmentsCreate.View {
	declarationNode := testCreateBasicDeclarationBooleanNode(name, "modifiable")

	b, _ := json.Marshal(value)
	handler := assignmentsCreate.New(assignmentsCreate.NewCreateNodeModel(declarationNode.Name, b))

	view, err := handler.Handle()
	testAssertErrNil(err)
	testAssertIDValid(view.ID)

	return view
}

func testAssertErrNil(err error) {
	gomega.Expect(err).Should(gomega.BeNil())
}

func testAssertIDValid(id string) {
	gomega.Expect(id).ShouldNot(gomega.BeEmpty())
	_, err := ulid.Parse(id)
	gomega.Expect(err).Should(gomega.BeNil())
}
