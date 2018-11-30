package projectify

import (
	"fmt"
	"strconv"
)

type StructNode struct {
	Connections []*StructNode
	value       string
	id          int
}

func (ref StructNode) New(id int, value string) StructNode {
	n := StructNode{[]*StructNode{}, value, id}
	return n
}

/*
*
 */
func (ref *StructNode) isConnected(node *StructNode) bool {
	if ref == node {
		return true
	} else {
		found := false
		// Check for node in ref.
		for i := 0; i < len(ref.Connections); i++ {
			if ref.Connections[i] == node {
				found = found || ref.Connections[i].isConnected(node)
			}
		}
		// Check for ref in node.
		for i := 0; i < len(node.Connections); i++ {
			if node.Connections[i] == ref {
				found = found || node.Connections[i].isConnected(ref)
			}
		}
		return found
	}
}

func (ref *StructNode) AddConnection(node *StructNode) bool {
	self := *ref
	if !ref.isConnected(node) {
		arr := append(self.Connections, node)
		ref.Connections = arr
		return true
	} else {
		return false
	}
}

func (ref StructNode) Print(i int) {
	fmt.Println(ref.value + ":" + strconv.Itoa(i))
	for i := 0; i < len(ref.Connections); i++ {
		ref.Connections[i].Print(i + 1)
	}
}
