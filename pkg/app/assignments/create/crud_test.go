package create

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/segmentio/ksuid"
)

var _ = ginkgo.Describe("Assignment CRUD success test", func() {
	ginkgo.It("should create an assignment text node when the node does not exists", ginkgo.Label("assignment", "crud", "success", "1"), func() {
		name := ksuid.New().String()
		declarationNode := testCreateBasicDeclarationTextNode(name, "modifiable")

		text := "this is a text node"
		b, _ := json.Marshal(text)

		handler := New(NewCreateNodeModel(declarationNode.Name, b))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())

		gomega.Expect(view.Value).Should(gomega.Equal(text))
	})

	ginkgo.It("should create an assignment boolean node when the node does not exists", ginkgo.Label("assignment", "crud", "success", "2"), func() {
		name := ksuid.New().String()
		declarationNode := testCreateBasicDeclarationTextNode(name, "modifiable")

		b, _ := json.Marshal(false)
		handler := New(NewCreateNodeModel(declarationNode.Name, b))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())

		gomega.Expect(view.Value).Should(gomega.BeFalse())
	})

	ginkgo.It("should update an assignment text node when the node already exists", ginkgo.Label("assignment", "crud", "success", "3"), func() {
		name := ksuid.New().String()
		testCreateBasicAssignmentTextNode(name)
		text := "this is a changed text value"

		b, _ := json.Marshal(text)

		handler := New(NewCreateNodeModel(name, b))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())

		gomega.Expect(view.Value).Should(gomega.Equal(text))
	})

	ginkgo.It("should update an assignment boolean node when the node already exists", ginkgo.Label("assignment", "crud", "success", "4"), func() {
		name := ksuid.New().String()
		testCreateBasicAssignmentBooleanNode(name, true)

		b, _ := json.Marshal(false)
		handler := New(NewCreateNodeModel(name, b))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())

		gomega.Expect(view.Value).Should(gomega.BeFalse())

	})

	ginkgo.It("should update an assignment node from text to boolean", ginkgo.Label("assignment", "crud", "success", "5"), func() {
		name := ksuid.New().String()
		testCreateBasicAssignmentTextNode(name)

		b, _ := json.Marshal(false)
		handler := New(NewCreateNodeModel(name, b))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())

		gomega.Expect(view.Value).Should(gomega.BeFalse())
	})

	ginkgo.It("should update an assignment node from boolean to text", ginkgo.Label("assignment", "crud", "success", "5"), func() {
		name := ksuid.New().String()
		testCreateBasicAssignmentBooleanNode(name, true)
		text := "this is a text value"

		b, _ := json.Marshal(text)
		handler := New(NewCreateNodeModel(name, b))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID.String())

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())

		gomega.Expect(view.Value).Should(gomega.Equal(text))
	})
})
