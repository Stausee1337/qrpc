package analysis

import (
	"fmt"
	"strings"

	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

type Type interface{
	IsIRType()
	GoTypeRepr() string
}
type UserDefinedType interface{
	Type
	IsStruct() bool
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
func (n *NumberType) GoTypeRepr() string {
	switch n.Kind {
	case SIGNED:
		return "int"
	case UNSIGNED:
		return "uint"
	case FLOATING:
		return "float64"
	default:
		panic("invalid number type")
	}
}

type AnyType struct { }
func (*AnyType) IsIRType() {}
func (*AnyType) GoTypeRepr() string { return "any" }

type BoolType struct { }
func (*BoolType) IsIRType() {}
func (*BoolType) GoTypeRepr() string { return "bool" }

type EmptyType struct { }
func (*EmptyType) IsIRType() {}
func (*EmptyType) GoTypeRepr() string { panic("empty type cannot be represented in go") }

type StringType struct { }
func (*StringType) IsIRType() {}
func (*StringType) GoTypeRepr() string { return "string" }

type UUIDType struct { }
func (*UUIDType) IsIRType() {}
func (*UUIDType) GoTypeRepr() string { return "uuid.UUID" }

type ArrayType struct {
	Type Type
}
func (*ArrayType) IsIRType() {}
func (a *ArrayType) GoTypeRepr() string {
	return fmt.Sprintf("[]%v", a.Type.GoTypeRepr())
}

type UnionType struct {
	Types []Type
}
func (*UnionType) IsIRType() {}
func (u *UnionType) GoTypeRepr() string {
	variants := make([]string, 0)
	for idx, ty := range u.Types {
		variants = append(
			variants,
			fmt.Sprintf("Variant%v *%v `json:\"variant%v\"`", idx + 1, ty.GoTypeRepr(), idx + 1),
		)
	}
	return fmt.Sprintf(
		"struct { %v }",
		strings.Join(variants, "; "),
	)
}

type OptionalType struct {
	Type Type
}
func (*OptionalType) IsIRType() {}
func (o *OptionalType) GoTypeRepr() string {
	return fmt.Sprintf("*%v", o.Type.GoTypeRepr())
}

type EnumType struct {
	Name     parser.Symbol
	Variants []parser.Symbol
}

func (*EnumType) IsIRType() {}
func (*EnumType) IsStruct() bool { return false; }
func (o *EnumType) GoTypeRepr() string {
	return source.ToPascalCase(o.Name.String())
}

type RecordType struct {
	Name   parser.Symbol
	Fields []NamedType
}

func (*RecordType) IsIRType() {}
func (*RecordType) IsStruct() bool { return true; }
func (o *RecordType) GoTypeRepr() string {
	return source.ToPascalCase(o.Name.String())
}

type NamedType struct {
	Name parser.Symbol
	Type Type
}

type Service struct {
	Name       parser.Symbol
	Operations []Operation
}

func (s *Service) GoName() string {
	return source.ToPascalCase(s.Name.String())
}

func (s *Service) HasQueries() bool {
	for _, op := range s.Operations {
		if op.Kind == parser.OperationQuery {
			return true
		}
	}
	return false
}

func (s *Service) HasMutations() bool {
	for _, op := range s.Operations {
		if op.Kind == parser.OperationMutation {
			return true
		}
	}
	return false
}

type Operation struct {
	Kind       parser.OperationKind
	Name       parser.Symbol
	InputTypes []NamedType
	ResultType Type
}

func (o *Operation) GoName() string {
	return source.ToPascalCase(o.Name.String())
}

func (o *Operation) IsQuery() bool {
	return o.Kind == parser.OperationQuery
}

func (o *Operation) GoInputs() string {
	inputs := make([]string, 0)
	for _, namedTy := range o.InputTypes {
		inputs = append(inputs, fmt.Sprintf("%v %v", namedTy.Name, namedTy.Type.GoTypeRepr()))
	}
	return strings.Join(inputs, ", ")
}

func (o *Operation) HasInputs() bool {
	return len(o.InputTypes) > 0
}

func (o *Operation) HasOutput() bool {
	_, isEmpty := o.ResultType.(*EmptyType)
	return !isEmpty
}

func (o *Operation) GoOutputs() string {
	_, isEmpty := o.ResultType.(*EmptyType)
	if isEmpty {
		return "error"
	}
	return fmt.Sprintf("%v, error", o.ResultType.GoTypeRepr())
}

func (o *Operation) RealName() string {
	return o.Name.String()
}

type AnalysisResult struct {
	Types     []UserDefinedType
	Services  []Service
}

