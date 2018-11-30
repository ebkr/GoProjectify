package projectify

type StructProject struct {
	templates map[*StructNode]string
	binds     map[*StructNode]*StructNode
}

func (ref *StructProject) AddTemplate(node *StructNode) {
	ref.templates[node] = node.value
}

func (ref *StructProject) BindNodes(nodeA *StructNode, nodeB *StructNode) {
	if nodeA.AddConnection(nodeB) {
		ref.binds[nodeA] = nodeB
	}
}
