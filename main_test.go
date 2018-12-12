package main

import (
	"log"
	"os"
	"strconv"
	"testing"

	"./ApplicationData/Libraries/projectify"
)

// generateDirectoryAndFile : Create a temporary directory, and a file inside it.
func generateDirectoryAndFile(test *testing.T) (*projectify.StructCreate, *projectify.StructCreate) {
	tempDir := projectify.StructCreate{}.New(".", "/Test")
	if !tempDir.CheckExistence() {
		os.Mkdir("./Test", os.FileMode(0755))
	}
	if !tempDir.CheckExistence() {
		test.Errorf("Could not create directory")
	}
	file := projectify.StructCreate{}.New("./Test/", "MyFile.projectify")
	file.OverwriteFile("<<TEMPLATES>>\n<<BINDS>>\n<<POSITIONS>>\n")
	if !file.CheckExistence() {
		test.Errorf("File not created")
	}
	return &tempDir, &file
}

// updateProjectTree : Updates nodes for the project
func updateProjectTree(project *projectify.StructProject, file *projectify.StructCreate) {
	log.Println("Updating Tree")
	log.Println(file.GetData())
	nodes := file.GenerateNodeTree()
	log.Println(nodes)
	myMap := map[*projectify.StructNode]string{}
	for i := 0; i < len(nodes); i++ {
		log.Println("Found Node: " + nodes[i].GetValue())
		myMap[nodes[i]] = nodes[i].GetValue()
	}
	log.Println("Setting Tree")
	project.SetTree(myMap)
}

// Test_CreateAndRemoveDirectory : Creates a directory, and attempts to remove it
func Test_CreateAndRemoveDirectory(test *testing.T) {
	log.Println("> Running Test: CreateAndRemoveDirectory")
	dir, file := generateDirectoryAndFile(test)
	if file.CheckExistence() {
		file.Delete()
	}
	if dir.CheckExistence() {
		dir.Delete()
		os.Remove("./Test")
	}
	if dir.CheckExistence() {
		test.Errorf("Directory still exists")
	}
}

// Test_CreateWriteAndRemoveFile : Creates a directory with a file inside, writes to the file, and deletes.
func Test_CreateWriteAndRemoveFile(test *testing.T) {
	log.Println("> Running Test: CreateWriteAndRemoveFile")
	tempDir, file := generateDirectoryAndFile(test)
	file.OverwriteFile("MyData")
	file.Delete()
	if file.CheckExistence() {
		test.Errorf("File could not be deleted")
	}
	tempDir.Delete()
}

// Test_NodeTest : Creates two nodes, and attempts to link them together. Should not allow recursive linking.
func Test_NodeTest(test *testing.T) {
	log.Println("> Running Test: NodeTest")
	node := projectify.StructNode{}.New(1, "Node Value", 0, 0)
	node2 := projectify.StructNode{}.New(2, "Other Node", 0, 0)
	if !node.AddConnection(&node2) {
		test.Errorf("Failed to join nodes")
	} else {
		if node2.AddConnection(&node) {
			test.Errorf("Node should not connect to parent")
		}
	}
}

// Test_GenerateSimpleNodeTree : Creates a demo project, and generates a specific node tree
func Test_GenerateSimpleNodeTree(test *testing.T) {
	log.Println("> Running Test: GenerateSimpleNodeTree")
	tempDir, _ := generateDirectoryAndFile(test)
	file := projectify.StructCreate{}.New("./Test", "/MyFile.projectify")
	project := projectify.StructProject{}
	project.Init()
	file.NewNode(0, "Example1")
	file.NewNode(1, "Example2")
	log.Println(">> Updating Project Tree #1")
	updateProjectTree(&project, &file)
	log.Println(">> Adding binds")
	nodeA := project.GetNodeByID(0)
	nodeB := project.GetNodeByID(1)
	if nodeA == nil {
		test.Errorf("Node A cannot be found")
	}
	if nodeB == nil {
		test.Errorf("Node B cannot be found")
	}
	if project.GetNodeByID(0).AddConnection(project.GetNodeByID(1)) {
		log.Println(">> Applying Bind")
		file.AppendFile("<<BINDS>>", "0:1")
	}
	log.Println(">> Updating Project Tree #2")
	updateProjectTree(&project, &file)
	log.Println(">> Getting Project Tree")
	tree := project.GetTree()
	var counter int
	for node := range tree {
		counter++
		log.Println(">> Loop")
		if node.GetID() == 0 {
			for _, node := range node.Connections {
				log.Println(node.GetValue())
			}
			if len(node.Connections) != 1 {
				test.Errorf("Wrong number of connections")
			}
		}
	}
	if counter != 2 {
		test.Errorf("Wrong number of nodes created")
	}
	log.Println(">> Delete")
	file.Delete()
	tempDir.Delete()
}

// Test_GenerateExtendedNodeTree : Creates a demo project, and generates a node tree with duplicate links
func Test_GenerateExtendedNodeTree(test *testing.T) {
	log.Println("> Running Test: GenerateExtendedNodeTree")
	tempDir, file := generateDirectoryAndFile(test)
	project := projectify.StructProject{}
	project.Init()
	file.AppendFile("<<TEMPLATES>>", "0:Example1")
	file.AppendFile("<<TEMPLATES>>", "1:Example2")
	updateProjectTree(&project, file)
	if project.GetNodeByID(0).AddConnection(project.GetNodeByID(1)) {
		file.AppendFile("<<BINDS>>", "0:1")
		file.AppendFile("<<BINDS>>", "0:1")
	}
	updateProjectTree(&project, file)
	tree := project.GetTree()
	var counter int
	for node := range tree {
		counter++
		if node.GetID() == 0 {
			for _, node := range node.Connections {
				log.Println(node.GetValue())
			}
			if len(node.Connections) != 1 {
				test.Errorf("Wrong number of connections")
			}
		}
	}
	if counter != 2 {
		test.Errorf("Wrong number of nodes created")
	}
	file.Delete()
	tempDir.Delete()
}

// Test_GetUniqueID : Attempts to retrieve the correct AvailableID
func Test_GetUniqueID(test *testing.T) {
	log.Println("> Running Test: GetUniqueID")
	// Create Defaults
	tempDir, file := generateDirectoryAndFile(test)
	project := projectify.StructProject{}

	// Initialise and update project
	updateProjectTree(&project, file)

	// Add Nodes
	file.NewNode(0, "Example1")
	file.NewNode(1, "Example2")
	file.RemoveNode(0)

	// Update project
	updateProjectTree(&project, file)

	// Attempt to get new NodeID
	var id int = project.GetAvailableID()
	if id != 0 {
		test.Errorf("Invalid ID #1. Expecting 0, got " + strconv.Itoa(id))
	}
	file.NewNode(project.GetAvailableID(), "Example1")

	// Update Again
	updateProjectTree(&project, file)

	// Get new ID
	id = project.GetAvailableID()
	if id != 2 {
		test.Errorf("Invalid ID #2. Expecting 2, got " + strconv.Itoa(id))
	}

	file.Delete()
	tempDir.Delete()
}
