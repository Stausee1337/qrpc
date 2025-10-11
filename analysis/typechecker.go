package analysis

import (
	"fmt"
	"slices"

	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

type typechecker struct {
	queryCache 		map[*parser.Item]UserDefinedType
	queryStack		[]*parser.Item
	primitiveTypes  map[parser.Symbol]Type
}

func makeTypechecker() typechecker {
	return typechecker{
		queryCache: map[*parser.Item]UserDefinedType{},
		queryStack: []*parser.Item{},
		primitiveTypes: map[parser.Symbol]Type{
			parser.SYMstring: &StringType{},
			parser.SYMany: &AnyType{},
			parser.SYMbool: &BoolType{},
			parser.SYMempty: &EmptyType{},
			parser.SYMfloat: &NumberType{ Kind: FLOATING },
			parser.SYMint: &NumberType{ Kind: SIGNED },
			parser.SYMuint: &NumberType{ Kind: UNSIGNED },
			parser.SYMuuid: &UUIDType{},
		},
	};
}

func (c *typechecker) analyze(items []parser.Item) (AnalysisResult, *source.SourceError) {
	types := make([]UserDefinedType, 0)
	services := make([]Service, 0)
	for i := range items {
		item := &items[i]

		switch v := item.Data.(type) {
		case *parser.IEnum, *parser.IRecord:
			ty, err := c.queryNode(item)
			if err != nil {
				return AnalysisResult{}, err
			}
			types = append(types, ty)
		case *parser.IService:
			service, err := c.buildService(v)
			if err != nil {
				return AnalysisResult{}, err
			}
			services = append(services, service)
		}
	}

	return AnalysisResult{
		Types: types,
		Services: services,
	}, nil
}

func (c *typechecker) queryRef(item *parser.TRef) (Type, *source.SourceError) {
	switch v := item.Res.(type) {
	case *parser.Item:
		return c.queryNode(v)
	case parser.Symbol:
		return c.primitiveTypes[v], nil
	default:
		panic("invalid resolution")
	}
}

func (c *typechecker) queryNode(item *parser.Item) (UserDefinedType, *source.SourceError) {
	if slices.Contains(c.queryStack, item) {
		return nil, &source.SourceError{
			Pos: item.Pos,
			Message: fmt.Sprintf("type '%v' cannot be recursive", item.Data.GetName().Symbol),
		}
	}
	udt, ok := c.queryCache[item]
	if ok {
		return udt, nil
	}

	c.queryStack = append(c.queryStack, item)
	udt, err := c.buildUDT(item)
	if err != nil {
		return nil, err
	}
	pop(&c.queryStack)
	c.queryCache[item] = udt

	return udt, nil 
}

func (c *typechecker) buildUDT(item *parser.Item) (UserDefinedType, *source.SourceError) {
	switch v := item.Data.(type) {
	case *parser.IEnum:
		return c.buildEnum(v);
	case *parser.IRecord:
		return c.buildRecord(v);
	default:
		panic("invalid type in buildUDT")
	}
}

func (c *typechecker) buildEnum(enum *parser.IEnum) (UserDefinedType, *source.SourceError) {
	variants := make([]parser.Symbol, 0)
	for _, variant := range enum.Variants {
		variants = append(variants, variant.Symbol)
	}

	return &EnumType {
		Name: enum.Name.Symbol,
		Variants: variants,
	}, nil;
}

func (c *typechecker) buildRecord(record *parser.IRecord) (UserDefinedType, *source.SourceError) {
	fields := make([]NamedType, 0)
	for idx := range record.Fields {
		field := &record.Fields[idx]
		ty, err := c.transformType(&field.Type)
		if err != nil {
			return nil, err;
		}

		fields = append(fields, NamedType{
			Name: field.Name.Symbol,
			Type: ty,
		})
	}

	return &RecordType {
		Name: record.Name.Symbol,
		Fields: fields,
	}, nil;
}

func (c *typechecker) buildService(service *parser.IService) (Service, *source.SourceError) {
	operations := make([]Operation, 0)
	for idx := range service.Operations {
		op, err := c.buildOperation(&service.Operations[idx])
		if err != nil {
			return Service{}, err;
		}

		operations = append(operations, op)
	}

	return Service {
		Name: service.Name.Symbol,
		Operations: operations,
	}, nil;
}

func (c *typechecker) buildOperation(op *parser.Operation) (Operation, *source.SourceError) {
	inputs := make([]NamedType, 0)
	for idx := range op.Inputs {
		input := &op.Inputs[idx]
		ty, err := c.transformType(&input.Type)
		if err != nil {
			return Operation{}, err;
		}

		inputs = append(inputs, NamedType{
			Name: input.Name.Symbol,
			Type: ty,
		})
	}
	result, err := c.transformType(&op.ResultType)
	if err != nil {
		return Operation{}, err;
	}

	return Operation {
		Kind: op.Kind,
		Name: op.Name.Symbol,
		InputTypes: inputs,
		ResultType: result,
	}, nil;
}

func (c *typechecker) transformType(ty *parser.Type) (Type, *source.SourceError) {
	switch v := ty.Data.(type) {
	case *parser.TRef:
		return c.transformTRef(v);
	case *parser.TArray:
		return c.transformTArray(v);
	case *parser.TUnion:
		return c.transformTUnion(v);
	case *parser.TOptional:
		return c.transformTOptional(v);
	default:
		panic("unknown type variant")
	}
}

func (c *typechecker) transformTRef(ref *parser.TRef) (Type, *source.SourceError) {
	return c.queryRef(ref)
}

func (c *typechecker) transformTArray(array *parser.TArray) (Type, *source.SourceError) {
	ty, err := c.transformType(&array.Base)
	if err != nil { return nil, err; }
	return &ArrayType { Type: ty }, nil;
}

func (c *typechecker) transformTUnion(union *parser.TUnion) (Type, *source.SourceError) {
	types := make([]Type, 0)
	for idx := range union.Variants {
		ty, err := c.transformType(&union.Variants[idx])
		if err != nil { return nil, err; }
		types = append(types, ty)
	}
	return &UnionType { Types: types }, nil;
}

func (c *typechecker) transformTOptional(optional *parser.TOptional) (Type, *source.SourceError) {
	ty, err := c.transformType(&optional.Base)
	if err != nil { return nil, err; }
	return &OptionalType { Type: ty }, nil;
}

func pop[T any](list *[]T) T {
    l := *list
    v := l[len(l)-1]
    *list = l[:len(l)-1]
    return v
}
