package main

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/ebkr/GoProjectify/ApplicationData/Libraries/projectify"
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
	project.Init()
	nodes := file.GenerateNodeTree()
	myMap := map[*projectify.StructNode]string{}
	for i := 0; i < len(nodes); i++ {
		myMap[nodes[i]] = nodes[i].GetValue()
	}
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

// Test_GetUniqueID : Attempts to retrieve the correct AvailableID
func Test_GetUniqueID(test *testing.T) {
	log.Println("> Running Test: GetUniqueID")

	project := projectify.StructProject{}
	project.Init()

	// Add Nodes

	nodeA := projectify.StructNode{}.New(0, "Example1", 0, 0)
	nodeB := projectify.StructNode{}.New(1, "Example2", 0, 0)

	tree := map[*projectify.StructNode]string{
		&nodeB: nodeB.GetValue(),
	}

	if !nodeA.AddConnection(&nodeB) {
		test.Errorf("Could not link nodes")
	}

	project.SetTree(tree)

	// Attempt to get new NodeID
	var id int = project.GetAvailableID()
	if id != 0 {
		test.Errorf("Invalid ID #1. Expecting 0, got " + strconv.Itoa(id))
	}

	tree = map[*projectify.StructNode]string{
		&nodeB: nodeB.GetValue(),
		&nodeA: nodeA.GetValue(),
	}

	project.SetTree(tree)

	// Get new ID
	id = project.GetAvailableID()
	if id != 2 {
		test.Errorf("Invalid ID #2. Expecting 2, got " + strconv.Itoa(id))
	}
}
