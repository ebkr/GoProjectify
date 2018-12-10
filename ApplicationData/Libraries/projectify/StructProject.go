/*

Package projectify::StructProject

Used to create a project object, that stores project nodes, and provide information about them

*/

package projectify

import (
	"sort"
)

// StructProject : Contains a list of all nodes associated with loaded project
type StructProject struct {
	templates map[*StructNode]string
}

// Init : Initialise empty map
func (ref *StructProject) Init() {
	ref.templates = make(map[*StructNode]string)
}

// SetTree : Overwrite templates map
func (ref *StructProject) SetTree(nodes map[*StructNode]string) {
	ref.templates = nodes
}

// GetAvailableID : Retrieve the lowest possible ID. Prevents ridiculous IDs.
func (ref *StructProject) GetAvailableID() int {
	id := -1
	numeric := []int{}
	for k := range ref.templates {
		numeric = append(numeric, k.GetID())
	}
	sort.Ints(numeric)
	for i := 0; i < len(numeric); i++ {
		if numeric[i] == id+1 {
			id++
		}
	}
	return (id + 1)
}

// GetNodeByName : Return a StructNode via value search
func (ref *StructProject) GetNodeByName(name string) *StructNode {
	for k, v := range ref.templates {
		if v == name {
			return k
		}
	}
	return nil
}

// GetNodeByID : Return a StructNode via an ID search
func (ref *StructProject) GetNodeByID(id int) *StructNode {
	for k := range ref.templates {
		if k.GetID() == id {
			return k
		}
	}
	return nil
}

// GetTree : Return the project's tree
func (ref *StructProject) GetTree() map[*StructNode]string {
	return ref.templates
}
