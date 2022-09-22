package app

import "errors"

var (
	ErrBackReferenceError = errors.New("back reference error")
	ErrEmptyInput         = errors.New("input is empty")
)
