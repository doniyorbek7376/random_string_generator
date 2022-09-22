package app

import (
	"strconv"
	"strings"

	"github.com/doniyorbek7376/random_string_generator/models"
)

const asciiCharacters = `!"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~`

// Parser creates generator tree based on tokens.
type ParserI interface {
	Parse(tokens []models.Token, rootNode RootNode) (Node, error)
	ResetGroupResults()
}

type parser struct {
	groupResults []string
}

func NewParser() ParserI {
	return &parser{make([]string, 0)}
}

func (p *parser) ResetGroupResults() {
	p.groupResults = make([]string, 0)
}

func (p *parser) Parse(tokens []models.Token, rootNode RootNode) (Node, error) {
	nodes := NewRootNode()
	n := len(tokens)
	for i := 0; i < n; i++ {
		token := tokens[i]
		switch token.Type {
		case models.ClassOpener:
			// When we find opening bracket, we locate the corresponding closing pair for it, and
			// process the tokens in between separately.
			j := i + 1
			for ; j < n; j++ {
				if tokens[j].Type == models.ClassCloser {
					break
				}
			}
			if j == n {
				// Could not find closing bracket pair
				return nil, ErrClassBracketNotClosed
			}
			node, err := p.parseClassTokens(tokens[i+1 : j])
			if err != nil {
				return nil, err
			}
			nodes.AddChild(node)
			i = j

		case models.GroupOpener:
			j := i + 1
			depth := 1
			for ; j < n; j++ {
				if tokens[j].Type == models.GroupOpener {
					depth++
				}
				if tokens[j].Type == models.GroupCloser {
					depth--
				}
				if depth == 0 {
					break
				}
			}
			if j == n {
				return nil, ErrGroupBracketNotClosed
			}

			node, err := p.Parse(tokens[i+1:j], NewGroupNode(&p.groupResults))
			if err != nil {
				return nil, err
			}
			nodes.AddChild(node)
			i = j

		case models.CounterOpener:
			j := i + 1
			for ; j < n; j++ {
				if tokens[j].Type == models.CounterCloser {
					break
				}
			}
			if j == n {
				return nil, ErrCounterBracketNotClosed
			}
			node, err := p.parseCounterTokens(tokens[i+1:j], nodes.PopChild())
			if err != nil {
				return nil, err
			}
			nodes.AddChild(node)
			i = j

		case models.ClassCloser, models.GroupCloser, models.CounterCloser:
			// Closing brackets without opening pairs.
			return nil, ErrInvalidClosingBracket

		case models.Asteriks:
			// Asteriks multiplies last node by [min,max] interval of [0, inf).
			// We pop the last node and wrap it with MultiplyNode
			lastNode := nodes.PopChild()
			if lastNode == nil {
				return nil, ErrInvalidAsteriks
			}
			// For simplicity, if the max value of the interval is infinite, I capped it with 10
			nodes.AddChild(NewMultiplyNode(lastNode, 0, 10))

		case models.QuestionMark:
			// Same as with the Asteriks, but the [min, max] interval is [0, 1]
			lastNode := nodes.PopChild()
			if lastNode == nil {
				return nil, ErrInvalidQuestionMark
			}
			nodes.AddChild(NewMultiplyNode(lastNode, 0, 1))

		case models.Plus:
			// Same as with the Asteriks, but [min-max] interval is [1, inf)
			lastNode := nodes.PopChild()
			if lastNode == nil {
				return nil, ErrInvalidPlus
			}
			nodes.AddChild(NewMultiplyNode(lastNode, 1, 100))

		case models.Dot:
			// Any ascii character
			nodes.AddChild(NewRandomNode(strings.Split(asciiCharacters, "")))

		case models.AlternatingBranch:
			// All the nodes accumulated so far, go to the left of the AlternateNode.
			// All the remaininng tokens after getting parsed become right child of the node.
			left := nodes
			right, err := p.Parse(tokens[i+1:], NewRootNode())
			if err != nil {
				return nil, err
			}
			// As previous root node became left child of the AlternateNode,
			// we need a new root node
			nodes = NewRootNode()
			nodes.AddChild(NewAlternateNode(left, right))
			i = n

		case models.BackSlash:
			// If the backslash is followed by a digit, it's a back reference.
			// Else, it's just a constant value
			if i == n-1 {
				return nil, ErrInvalidBackSlash
			}
			i++
			nextToken := tokens[i]
			groupIndex, err := strconv.Atoi(nextToken.Value)
			if err != nil {
				// Could not convert to int, so treating it as a literal value
				nodes.AddChild(NewTextNode(nextToken.Value))
				continue
			}
			if groupIndex == 0 {
				// Back reference should start with 1
				return nil, ErrBackReferenceError
			}
			groupIndex--
			nodes.AddChild(NewBackReferenceNode(groupIndex, &p.groupResults))

		default:
			nodes.AddChild(NewTextNode(token.Value))
		}
	}

	rootNode.AddChild(nodes)
	return rootNode, nil
}

// Parses tokens inside character class brackets []
func (p *parser) parseClassTokens(tokens []models.Token) (Node, error) {
	var (
		negated bool
	)
	n := len(tokens)
	if n == 0 {
		return nil, ErrEmptyClass
	}
	valueSet := []string{}
	for i, token := range tokens {
		if token.Type == models.ClassNegater {
			negated = true
			continue
		}
		if token.Type == models.ClassRange && i != 0 && i != n-1 {
			nextToken := tokens[i+1]
			start := valueSet[len(valueSet)-1]
			end := nextToken.Value
			if start > end {
				return nil, ErrInvalidClassRange
			}
			for i := start[0] + 1; i < end[0]; i++ {
				valueSet = append(valueSet, string(i))
			}
			continue
		}
		valueSet = append(valueSet, token.Value)
	}
	newValueSet := []string{}
	m := toBoolMap(valueSet)
	for _, ch := range asciiCharacters {
		// if negated and the char is in map, or not negated and the char is not in map,
		// then we don't add that char to the result
		if negated && m[string(ch)] || !negated && !m[string(ch)] {
			continue
		}
		newValueSet = append(newValueSet, string(ch))
	}
	valueSet = newValueSet

	return NewRandomNode(valueSet), nil
}

func (p *parser) parseCounterTokens(tokens []models.Token, child Node) (Node, error) {
	var (
		min, max       int
		minStr, maxStr string
		err            error
		commaFound     bool
	)
	if child == nil {
		return nil, ErrInvalidCounter
	}
	for _, token := range tokens {
		if token.Type == models.Comma {
			if commaFound {
				return nil, ErrInvalidCounter
			}
			commaFound = true
			continue
		}
		if commaFound {
			maxStr += token.Value
			continue
		}
		minStr += token.Value
	}

	if maxStr != "" {
		max, err = strconv.Atoi(maxStr)
		if err != nil {
			return nil, err
		}
	}
	if minStr != "" {
		min, err = strconv.Atoi(minStr)
		if err != nil {
			return nil, err
		}
	}
	if max == 0 {
		max = min
	}

	return NewMultiplyNode(child, int(min), int(max)), nil
}

// helper function to create "exists" map from an array
func toBoolMap(arr []string) map[string]bool {
	res := make(map[string]bool)
	for _, val := range arr {
		res[val] = true
	}

	return res
}
