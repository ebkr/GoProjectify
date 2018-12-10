package main

import (
	"os"
	"testing"

	"./ApplicationData/Libraries/projectify"
)

// Test_CreateAndRemoveDirectory : Creates a directory, and attempts to remove it
func Test_CreateAndRemoveDirectory(test *testing.T) {
	tempDir := projectify.StructCreate{}.New("./Test", "")
	if !tempDir.CheckExistence() {
		os.Mkdir("./Test", os.FileMode(0755))
	}
	if !tempDir.CheckExistence() {
		test.Errorf("Could not create directory")
	} else {
		tempDir.Delete()
		if tempDir.CheckExistence() {
			test.Errorf("Could not delete directory")
		}
	}
}

// Test_CreateWriteAndRemoveFile : Creates a directory with a file inside, writes to the file, and deletes.
func Test_CreateWriteAndRemoveFile(test *testing.T) {
	tempDir := projectify.StructCreate{}.New("./Test", "")
	if !tempDir.CheckExistence() {
		os.Mkdir("./Test", os.FileMode(0755))
	}
	file := projectify.StructCreate{}.New("./Test/", "MyFile")
	file.OverwriteFile("MyData")
	file.Delete()
	if file.CheckExistence() {
		test.Errorf("File could not be deleted")
	}
	tempDir.Delete()
}

// Test_NodeTest : Creates two nodes, and attempts to link them together. Should not allow recursive linking.
func Test_NodeTest(test *testing.T) {
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
