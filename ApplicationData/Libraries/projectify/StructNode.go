package projectify

import (
	"fmt"
	"strconv"
)

// StructNode contains a value, id, and a list of all connected nodes.
// Nodes cannot be recursive.
type StructNode struct {
	Connections []*StructNode
	value       string
	id          int
	posX        float64
	posY        float64
}

// New : Create basic object
func (ref StructNode) New(id int, value string, x, y float64) StructNode {
	n := StructNode{[]*StructNode{}, value, id, x, y}
	return n
}

// isConnected : Compares connected nodes, to prevent recursive connections
func (ref *StructNode) isConnected(node *StructNode) bool {
	return node.recursiveFind(ref)
}

// recursiveFind : Attempts to find a node object in all connections, and subsequent connections.
func (ref *StructNode) recursiveFind(node *StructNode) bool {
	if ref == node {
		return true
	}
	for _, v := range ref.Connections {
		if v.recursiveFind(node) {
			return true
		}
	}
	return false
}

// AddConnection : Calls upon isConnected to determine if connection is valid. Will connect if valid.
func (ref *StructNode) AddConnection(node *StructNode) bool {
	if !ref.isConnected(node) {
		arr := append(ref.Connections, node)
		ref.Connections = arr
		return true
	}
	return false
}

// Print : Displays connected nodes
func (ref *StructNode) Print(i int) {
	fmt.Println(ref.value + ":" + strconv.Itoa(i))
	for i := 0; i < len(ref.Connections); i++ {
		ref.Connections[i].Print(i + 1)
	}
}

// GetID : Retrieve Node's ID
func (ref *StructNode) GetID() int {
	return ref.id
}

// GetValue : Retrieve Node's value
func (ref *StructNode) GetValue() string {
	return ref.value
}

// GetPosition : Return an array containing X:Y coordinates
func (ref *StructNode) GetPosition() [2]float64 {
	return [2]float64{ref.posX, ref.posY}
}

// SetPosition : Set X and Y coordinates
func (ref *StructNode) SetPosition(x, y float64) {
	ref.posX = x
	ref.posY = y
}
