package app

import (
	"math/rand"
	"strings"
)

// TextNode contains just a literal value
func (node *textNode) Generate() (string, error) {
	return node.value, nil
}

// RandomNode contains a slice of possible return values.
// Generator returns one random value from the set
func (node *randomNode) Generate() (string, error) {
	n := len(node.valueSet)
	if n == 0 {
		return "", nil
	}
	pos := rand.Intn(n)
	return node.valueSet[pos], nil
}

// MultiplyNode contains a child and [min, max] interval for number of generations.
// Generator returns concatenation of rand(min, max) generated values of its child.
func (node *multiplyNode) Generate() (string, error) {
	var sb strings.Builder
	count := rand.Intn(node.maxQuantity-node.minQuantity+1) + node.minQuantity
	for i := 0; i < count; i++ {
		childText, err := node.child.Generate()
		if err != nil {
			return "", err
		}
		sb.WriteString(childText)
	}
	return sb.String(), nil
}

// AlternateNode is for handling alternate branching (|).
// Returns the value of one of its two children.
func (node *alternateNode) Generate() (string, error) {
	n := rand.Intn(2)
	if n == 0 {
		return node.left.Generate()
	}
	return node.right.Generate()
}

// RootNode is a container for multiple nodes.
// Generates concatenation of values of its children.
func (node *rootNode) Generate() (string, error) {
	var sb strings.Builder
	for _, child := range node.children {
		childText, err := child.Generate()
		if err != nil {
			return "", err
		}
		sb.WriteString(childText)
	}
	return sb.String(), nil
}

// GroupNode is just like the RootNode, but saves the result in groupResults slice
// for back reference
func (node *groupNode) Generate() (string, error) {
	result, err := node.rootNode.Generate()
	if err != nil {
		return "", err
	}
	if node.groupResults != nil {
		*node.groupResults = append(*node.groupResults, result)
	}
	return result, nil
}

func (node *backReferenceNode) Generate() (string, error) {
	if len(*node.groupResults) <= node.index {
		return "", ErrBackReferenceError
	}
	return (*node.groupResults)[node.index], nil
}
