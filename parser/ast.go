package parser

import (
	"fmt"

	"github.com/stausee1337/qrpc/lexer"
	"github.com/stausee1337/qrpc/source"
)

type Ident struct {
	Symbol 	Symbol
	Pos 	source.Position
}

func makeIdent(tok *lexer.Token) Ident {
	if tok.Kind != lexer.Ident {
		panic(fmt.Sprintf("cannot make identifier from %v", tok.Kind.String()))
	}

	return Ident{
		Symbol: InternSymbol(tok.Value),
		Pos: tok.Pos,
	}
}

type Type struct {
	Pos	 source.Position
	Data T
}

type T interface { isType() }

func (t *TRef) isType() {}
func (t *TArray) isType() {}
func (t *TUnion) isType() {}
func (t *TOptional) isType() {}

type TRef struct {
	Name Ident
	Res  any
}

type TArray struct {
	Base Type
}

type TUnion struct {
	Variants []Type
}

type TOptional struct {
	Base Type
}

type Item struct {
	Pos	 source.Position
	Data I
}

type I interface {
	GetName() *Ident
}

func (i *IEnum) GetName() *Ident {
	return &i.Name;
}
func (i *IRecord) GetName() *Ident {
	return &i.Name;
}
func (i *IService) GetName() *Ident {
	return &i.Name;
}

type IEnum struct {
	Name 	 Ident
	Variants []Ident
}

type IRecord struct {
	Name 	Ident
	Fields 	[]Field
}

type Field struct {
	Name Ident
	Type Type
}

type IService struct {
	Name 		Ident
	Operations 	[]Operation
}

type OperationKind int

const (
	OperationQuery OperationKind = iota
	OperationMutation
)

type Operation struct {
	Pos			source.Position
	Kind 		OperationKind
	Name 		Ident
	Inputs		[]NamedInput
	ResultType	Type
}

type NamedInput struct {
	Pos	 source.Position
	Name Ident
	Type Type
}

