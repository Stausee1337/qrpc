package analysis

import (
	"github.com/stausee1337/qrpc/parser"
)

type Type interface{ IsIRType() }
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
func (*NumberType) IsIRType() {}

type AnyType struct { }
func (*AnyType) IsIRType() {}

type BoolType struct { }
func (*BoolType) IsIRType() {}

type EmptyType struct { }
func (*EmptyType) IsIRType() {}

type StringType struct { }
func (*StringType) IsIRType() {}

type UUIDType struct { }
func (*UUIDType) IsIRType() {}

type ArrayType struct {
	Type Type
}
func (*ArrayType) IsIRType() {}

type UnionType struct {
	Types []Type
}
func (*UnionType) IsIRType() {}

type OptionalType struct {
	Type Type
}
func (*OptionalType) IsIRType() {}

type EnumType struct {
	Name     parser.Symbol
	Variants []parser.Symbol
}

func (*EnumType) IsIRType() {}
func (*EnumType) isUDType() {}

type RecordType struct {
	Name   parser.Symbol
	Fields []NamedType
}

func (*RecordType) IsIRType() {}
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

