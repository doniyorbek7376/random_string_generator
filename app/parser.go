package app

import "github.com/doniyorbek7376/random_string_generator/models"

type ParserI interface {
	Parse([]models.Token) (Node, error)
}

type parser struct{}

func NewParser() ParserI {
	return parser{}
}

func (p parser) Parse([]models.Token) (Node, error) {
	panic("not implemented")
}
