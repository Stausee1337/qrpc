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
		return "Invalid"
	case EOF:
		return "EOF"
	case LParen:
		return "LParen"
	case RParen:
		return "RParen"
	case LBracket:
		return "LBracket"
	case RBracket:
		return "RBracket"
	case LBrace:
		return "LBrace"
	case RBrace:
		return "RBrace"
	case Colon:
		return "Colon"
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"
	case Question:
		return "Question"
	case VBar:
		return "VBar"
	case Ident:
		return "Ident"
	case KeywordEnum:
		return "KeywordEnum"
	case KeywordRecord:
		return "KeywordRecord"
	case KeywordService:
		return "KeywordService"
	case KeywordQuery:
		return "KeywordQuery"
	case KeywordMutation:
		return "KeywordMutation"
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
