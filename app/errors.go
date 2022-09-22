package app

import "errors"

var (
	ErrBackReferenceError      = errors.New("back reference error")
	ErrEmptyInput              = errors.New("input is empty")
	ErrClassBracketNotClosed   = errors.New("class bracket not closed, missing ]")
	ErrGroupBracketNotClosed   = errors.New("group bracket not closed, missing )")
	ErrCounterBracketNotClosed = errors.New("counter bracket not closed, missing }")
	ErrEmptyClass              = errors.New("class is empty")
	ErrInvalidClassRange       = errors.New("class range is invalid")
	ErrInvalidCounter          = errors.New("min max block invalid")
	ErrInvalidClosingBracket   = errors.New("invalid closing bracket")
	ErrInvalidAsteriks         = errors.New("invalid asteriks")
	ErrInvalidQuestionMark     = errors.New("invalid question mark")
	ErrInvalidPlus             = errors.New("invalid plus")
	ErrInvalidBackSlash        = errors.New("invalid backslash")
)
