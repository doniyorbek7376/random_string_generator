package app

import (
	"strings"

	"github.com/doniyorbek7376/random_string_generator/models"
)

type TokenizerI interface {
	Tokenize(input string) ([]models.Token, error)
}

type tokenizer struct {
}

func NewTokenizer() TokenizerI {
	return tokenizer{}
}

func (tk tokenizer) Tokenize(input string) ([]models.Token, error) {
	var (
		res          []models.Token
		classStarted bool
		backSlashed  bool
	)

	if len(input) == 0 {
		return nil, ErrEmptyInput
	}
	for _, char := range strings.Split(input, "") {
		token := models.Token{Value: char, Type: models.ConstantValue}
		if backSlashed {
			backSlashed = false
			res = append(res, token)
			continue
		}
		if classStarted {
			switch char {
			case "]":
				classStarted = false
				token.Type = models.ClassCloser
			case "-":
				token.Type = models.ClassRange
			case "^":
				token.Type = models.ClassNegater
			}
		} else {
			switch char {
			case "[":
				classStarted = true
				token.Type = models.ClassOpener
			case "*":
				token.Type = models.Asteriks
			case ".":
				token.Type = models.Dot
			case "?":
				token.Type = models.QuestionMark
			case "\\":
				backSlashed = true
				token.Type = models.BackSlash
			case "+":
				token.Type = models.Plus
			case "|":
				token.Type = models.AlternatingBranch
			case ",":
				token.Type = models.Comma
			case "{":
				token.Type = models.CounterOpener
			case "}":
				token.Type = models.CounterCloser
			case "(":
				token.Type = models.GroupOpener
			case ")":
				token.Type = models.GroupCloser
			}
		}
		res = append(res, token)
	}
	return res, nil
}
