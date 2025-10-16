package lexer

import (
	"fmt"

	"github.com/stausee1337/qrpc/source"
)

const (
	Invalid Kind = iota
	EOF
	LParen
	RParen
	LBracket
	RBracket
	LBrace
	RBrace
	Colon
	Comma
	Semicolon
	Question
	VBar

	Ident

	KeywordEnum
	KeywordRecord
	KeywordService
	KeywordQuery
	KeywordMutation
)

func (k Kind) String() string {
	switch k {
	case Invalid:
		return "<error>"
	case EOF:
		return "<eof>"
	case LParen:
		return "("
	case RParen:
		return ")"
	case LBracket:
		return "["
	case RBracket:
		return "]"
	case LBrace:
		return "{"
	case RBrace:
		return "}"
	case Colon:
		return ":"
	case Comma:
		return ","
	case Semicolon:
		return ";"
	case Question:
		return "?"
	case VBar:
		return "|"
	case Ident:
		return "<name>"
	case KeywordEnum:
		return "enum"
	case KeywordRecord:
		return "record"
	case KeywordService:
		return "service"
	case KeywordQuery:
		return "query"
	case KeywordMutation:
		return "mutation"
	}
	panic("unreachable")
}

type Kind int

type Token struct {
	Pos   	source.Position
	Kind  	Kind
	Value 	string
}

func (t *Token) String() string {
	if t.Value == "" {
		return fmt.Sprintf("%v@%v:%v", t.Kind.String(), t.Pos.Lineno, t.Pos.Column);
	}
	return fmt.Sprintf("%v@%v:%v \"%v\"", t.Kind.String(), t.Pos.Lineno, t.Pos.Column, t.Value);
}
