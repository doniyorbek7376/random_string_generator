package app

import "github.com/doniyorbek7376/random_string_generator/models"

type TokenizerI interface {
	Tokenize(input string) ([]models.Token, error)
}

type tokenizer struct {
}

func NewTokenizer() TokenizerI {
	return tokenizer{}
}

func (tk tokenizer) Tokenize(input string) ([]models.Token, error) {
	panic("not implemented")
}
