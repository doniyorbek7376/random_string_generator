package app

import (
	"math/rand"
	"strings"
)

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

func (node *textNode) Generate() (string, error) {
	return node.value, nil
}

func (node *randomNode) Generate() (string, error) {
	n := len(node.valueSet)
	if n == 0 {
		return "", nil
	}
	pos := rand.Intn(n)
	return node.valueSet[pos], nil
}

func (node *multiplyNode) Generate() (string, error) {
	var sb strings.Builder
	count := rand.Intn(node.maxQuantity-node.minQuantity) + node.minQuantity
	for i := 0; i < count; i++ {
		childText, err := node.child.Generate()
		if err != nil {
			return "", err
		}
		sb.WriteString(childText)
	}
	return sb.String(), nil
}

func (node *alternateNode) Generate() (string, error) {
	n := rand.Intn(2)
	if n == 0 {
		return node.left.Generate()
	}
	return node.right.Generate()
}

func (node *backReferenceNode) Generate() (string, error) {
	if len(*node.groupResults) <= node.index {
		return "", ErrBackReferenceError
	}
	return (*node.groupResults)[node.index], nil
}
