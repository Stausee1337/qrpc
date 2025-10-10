package analysis

import (
	"github.com/stausee1337/qrpc/parser"
)

type Type interface{ isIRType() }
type UserDefinedType interface{ isUDType() }

type EnumType struct {
	Name     parser.Symbol
	Variants []parser.Symbol
}

type RecordType struct {
	Name   parser.Symbol
	Fields []NamedType
}

type NamedType struct {
	Name parser.Symbol
	Type Type
}

type Service struct {
	Name       parser.Symbol
	Operations []Operation
}

type Operation struct {
	Kind       parser.OperationKind
	Name       parser.Symbol
	InputTypes []NamedType
	ResultType Type
}

type AnalysisResult struct {
	types     []UserDefinedType
	serivices []Service
}
