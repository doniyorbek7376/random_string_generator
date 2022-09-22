package models

type TokenType int

const (
	ConstantValue TokenType = iota
	Asteriks
	Dot
	Plus
	QuestionMark
	ClassOpener
	ClassCloser
	ClassNegater
	ClassRange
	GroupOpener
	GroupCloser
	AlternatingBranch
	BackReference
)
