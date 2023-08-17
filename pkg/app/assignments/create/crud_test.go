package create

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Assignment CRUD success test", func() {
	ginkgo.It("should create an assignment text node when the node does not exists", ginkgo.Label("assignment", "crud", "success", "1"), func() {
		name := uuid.NewString()
		declarationNode := testCreateBasicDeclarationTextNode(name, "modifiable")

		text := "this is a text node"
		b, _ := json.Marshal(text)

		handler := New(NewCreateNodeModel(declarationNode.Name, b))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())
		gomega.Expect(assignmentNode.ValueType).Should(gomega.Equal(assignments.ValueTextType))

		var textNode assignments.NodeText
		gomega.Expect(storage.Gorm().Where("assignment_node_id = ?", assignmentNode.ID).First(&textNode).Error).Should(gomega.BeNil())

		var createdText string
		gomega.Expect(json.Unmarshal(textNode.Value, &createdText))
		gomega.Expect(createdText).Should(gomega.Equal(text))
	})

	ginkgo.It("should create an assignment boolean node when the node does not exists", ginkgo.Label("assignment", "crud", "success", "1"), func() {
		name := uuid.NewString()
		declarationNode := testCreateBasicDeclarationTextNode(name, "modifiable")

		handler := New(NewCreateNodeModel(declarationNode.Name, true))

		view, err := handler.Handle()
		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())
		gomega.Expect(assignmentNode.ValueType).Should(gomega.Equal(assignments.ValueBooleanType))

		var booleanNode assignments.NodeBoolean
		gomega.Expect(storage.Gorm().Where("assignment_node_id = ?", assignmentNode.ID).First(&booleanNode).Error).Should(gomega.BeNil())

		gomega.Expect(booleanNode.Value).Should(gomega.BeTrue())
	})

	ginkgo.It("should update an assignment text node when the node already exists", ginkgo.Label("assignment", "crud", "success", "2"), func() {
		name := uuid.NewString()
		testCreateBasicAssignmentTextNode(name)
		text := "this is a changed text value"

		b, _ := json.Marshal(text)

		handler := New(NewCreateNodeModel(name, b))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())
		gomega.Expect(assignmentNode.ValueType).Should(gomega.Equal(assignments.ValueTextType))

		var textNode assignments.NodeText
		gomega.Expect(storage.Gorm().Where("assignment_node_id = ?", assignmentNode.ID).First(&textNode).Error).Should(gomega.BeNil())

		var changedText string
		gomega.Expect(json.Unmarshal(textNode.Value, &changedText))
		gomega.Expect(changedText).Should(gomega.Equal(text))
	})

	ginkgo.It("should update an assignment boolean node when the node already exists", ginkgo.Label("assignment", "crud", "success", "2"), func() {
		name := uuid.NewString()
		testCreateBasicAssignmentBooleanNode(name, true)

		handler := New(NewCreateNodeModel(name, false))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())
		gomega.Expect(assignmentNode.ValueType).Should(gomega.Equal(assignments.ValueBooleanType))

		var booleanNode assignments.NodeBoolean
		gomega.Expect(storage.Gorm().Where("assignment_node_id = ?", assignmentNode.ID).First(&booleanNode).Error).Should(gomega.BeNil())

		gomega.Expect(booleanNode.Value).Should(gomega.BeFalse())
	})

	ginkgo.It("should update an assignment node from text to boolean", ginkgo.Label("assignment", "crud", "success", "3"), func() {
		name := uuid.NewString()
		testCreateBasicAssignmentTextNode(name)

		handler := New(NewCreateNodeModel(name, true))

		view, err := handler.Handle()

		testAssertErrNil(err)
		testAssertIDValid(view.ID)

		var node declarations.Node
		gomega.Expect(storage.Gorm().Where("id = ?", view.ID).First(&node).Error).Should(gomega.BeNil())
		var assignmentNode assignments.Node
		gomega.Expect(storage.Gorm().Where("declaration_node_id = ?", node.ID).First(&assignmentNode).Error).Should(gomega.BeNil())
		gomega.Expect(assignmentNode.ValueType).Should(gomega.Equal(assignments.ValueBooleanType))

		var textNode assignments.NodeText
		gomega.Expect(storage.Gorm().Where("assignment_node_id = ?", assignmentNode.ID).First(&textNode).Error).ShouldNot(gomega.BeNil())

		var booleanNode assignments.NodeBoolean
		gomega.Expect(storage.Gorm().Where("assignment_node_id = ?", assignmentNode.ID).First(&booleanNode).Error).Should(gomega.BeNil())

		gomega.Expect(booleanNode.Value).Should(gomega.BeTrue())
	})
})
