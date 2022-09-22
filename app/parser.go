package app

import (
	"strconv"
	"strings"

	"github.com/doniyorbek7376/random_string_generator/models"
)

const asciiCharacters = `!"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~`

type ParserI interface {
	Parse(tokens []models.Token, rootNode RootNode) (Node, error)
}

type parser struct {
	groupResults []string
}

func NewParser() ParserI {
	return &parser{make([]string, 0)}
}

func (p *parser) Parse(tokens []models.Token, rootNode RootNode) (Node, error) {
	nodes := NewRootNode()
	n := len(tokens)
	for i := 0; i < n; i++ {
		token := tokens[i]
		switch token.Type {
		case models.ClassOpener:
			j := i + 1
			for ; j < n; j++ {
				if tokens[j].Type == models.ClassCloser {
					break
				}
			}
			if j == n {
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
			return nil, ErrInvalidClosingBracket

		case models.Asteriks:
			lastNode := nodes.PopChild()
			if lastNode == nil {
				return nil, ErrInvalidAsteriks
			}
			nodes.AddChild(NewMultiplyNode(lastNode, 0, 10))

		case models.QuestionMark:
			lastNode := nodes.PopChild()
			if lastNode == nil {
				return nil, ErrInvalidQuestionMark
			}
			nodes.AddChild(NewMultiplyNode(lastNode, 0, 1))

		case models.Plus:
			lastNode := nodes.PopChild()
			if lastNode == nil {
				return nil, ErrInvalidPlus
			}
			nodes.AddChild(NewMultiplyNode(lastNode, 1, 100))

		case models.Dot:
			nodes.AddChild(NewRandomNode(strings.Split(asciiCharacters, "")))

		case models.AlternatingBranch:
			left := nodes
			right, err := p.Parse(tokens[i+1:], NewRootNode())
			if err != nil {
				return nil, err
			}
			nodes = NewRootNode()
			nodes.AddChild(NewAlternateNode(left, right))
			i = n
		default:
			nodes.AddChild(NewTextNode(token.Value))
		}
	}

	rootNode.AddChild(nodes)
	return rootNode, nil
}

func (p *parser) parseClassTokens(tokens []models.Token) (Node, error) {
	var (
		negated bool
	)
	n := len(tokens)
	if n == 0 {
		return nil, ErrEmptyClass
	}
	valueSet := []string{}
	if negated && n == 1 {
		return nil, ErrEmptyClass
	}
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
			for i := start[0] + 1; i <= end[0]; i++ {
				valueSet = append(valueSet, string(i))
			}
			continue
		}
		valueSet = append(valueSet, token.Value)
	}
	if negated {
		newValueSet := []string{}
		negatedSet := toBoolMap(valueSet)
		for _, ch := range asciiCharacters {
			if negatedSet[string(ch)] {
				continue
			}
			newValueSet = append(newValueSet, string(ch))
		}
		valueSet = newValueSet
	}

	return NewRandomNode(valueSet), nil
}

func (p *parser) parseCounterTokens(tokens []models.Token, child Node) (Node, error) {
	var (
		min, max       int64
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
		max, err = strconv.ParseInt(maxStr, 10, 32)
		if err != nil {
			return nil, err
		}
	}
	if minStr != "" {
		min, err = strconv.ParseInt(minStr, 10, 32)
		if err != nil {
			return nil, err
		}
	}
	if max == 0 {
		max = min
	}

	return NewMultiplyNode(child, int(min), int(max)), nil
}

func toBoolMap(arr []string) map[string]bool {
	res := make(map[string]bool)
	for _, val := range arr {
		res[val] = true
	}

	return res
}
