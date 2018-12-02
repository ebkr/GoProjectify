/*

Package projectify::StructProject

Used to create a project object, that stores project nodes, and provide information about them

*/

package projectify

import (
	"sort"
)

type StructProject struct {
	templates map[*StructNode]string
}

func (ref *StructProject) Init() {
	ref.templates = make(map[*StructNode]string)
}

func (ref *StructProject) SetTree(nodes map[*StructNode]string) {
	ref.templates = nodes
}

func (ref *StructProject) GetAvailableId() int {
	id := -1
	numeric := []int{}
	for k := range ref.templates {
		numeric = append(numeric, k.GetId())
	}
	sort.Ints(numeric)
	for i := 0; i < len(numeric); i++ {
		if numeric[i] == id+1 {
			id++
		}
	}
	return (id + 1)
}

func (ref *StructProject) GetNodeByName(name string) *StructNode {
	for k, v := range ref.templates {
		if v == name {
			return k
		}
	}
	return nil
}

func (ref *StructProject) GetTree() map[*StructNode]string {
	return ref.templates
}
