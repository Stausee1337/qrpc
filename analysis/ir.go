package analysis

import (
	"github.com/stausee1337/qrpc/parser"
)

type Type interface{ isIRType() }
type UserDefinedType interface{
	Type
	isUDType()
}

const (
	SIGNED   NumberKind = "signed"
	UNSIGNED NumberKind = "unsigned"
	FLOATING NumberKind = "floating"
)

type NumberKind string

type NumberType struct {
	Kind NumberKind
}
func (*NumberType) isIRType() {}

type AnyType struct { }
func (*AnyType) isIRType() {}

type BoolType struct { }
func (*BoolType) isIRType() {}

type EmptyType struct { }
func (*EmptyType) isIRType() {}

type StringType struct { }
func (*StringType) isIRType() {}

type UUIDType struct { }
func (*UUIDType) isIRType() {}

type ArrayType struct {
	Type Type
}
func (*ArrayType) isIRType() {}

type UnionType struct {
	Types []Type
}
func (*UnionType) isIRType() {}

type OptionalType struct {
	Type Type
}
func (*OptionalType) isIRType() {}

type EnumType struct {
	Name     parser.Symbol
	Variants []parser.Symbol
}

func (*EnumType) isIRType() {}
func (*EnumType) isUDType() {}

type RecordType struct {
	Name   parser.Symbol
	Fields []NamedType
}

func (*RecordType) isIRType() {}
func (*RecordType) isUDType() {}

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
	Types     []UserDefinedType
	Services  []Service
}

