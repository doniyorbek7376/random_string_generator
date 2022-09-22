package app

type Node interface {
	Generate() (string, error)
}

type RootNode interface {
	AddChild(child Node)
}

type rootNode struct {
	children []Node
}

type groupNode struct {
	rootNode
	groupResults *[]string
}

type textNode struct {
	value string
}

type randomNode struct {
	valueSet []string
}

type multiplyNode struct {
	child       Node
	minQuantity int
	maxQuantity int
}

type alternateNode struct {
	left  Node
	right Node
}

type backReferenceNode struct {
	index        int
	groupResults *[]string
}

func (node *rootNode) AddChild(child Node) {
	node.children = append(node.children, child)
}
