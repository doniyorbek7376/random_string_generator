package app

type Node interface {
	Generate() (string, error)
}

type RootNode interface {
	Node
	AddChild(child Node)
}

type rootNode struct {
	children []Node
}

func (node *rootNode) AddChild(child Node) {
	node.children = append(node.children, child)
}

func NewRootNode() RootNode {
	return &rootNode{}
}

type groupNode struct {
	rootNode
	groupResults *[]string
}

func NewGroupNode(groupResults *[]string) RootNode {
	return &groupNode{groupResults: groupResults}
}

type textNode struct {
	value string
}

func NewTextNode(value string) Node {
	return &textNode{value}
}

type randomNode struct {
	valueSet []string
}

func NewRandomNode(valueSet []string) Node {
	return &randomNode{valueSet}
}

type multiplyNode struct {
	child       Node
	minQuantity int
	maxQuantity int
}

func NewMultiplyNode(child Node, min, max int) Node {
	return &multiplyNode{child, min, max}
}

type alternateNode struct {
	left  Node
	right Node
}

func NewAlternateNode(left, right Node) Node {
	return &alternateNode{left, right}
}

type backReferenceNode struct {
	index        int
	groupResults *[]string
}

func NewBackReferenceNode(index int, groupResults *[]string) Node {
	return &backReferenceNode{index, groupResults}
}
