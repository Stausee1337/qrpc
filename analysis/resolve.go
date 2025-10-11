package analysis

import (
	"fmt"

	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

type resolver struct {
	types map[parser.Symbol]*parser.Item
}

func makeResolver() resolver {
	return resolver{
		types: map[parser.Symbol]*parser.Item{},
	}
}

func (r *resolver) resolve(items []parser.Item) *source.SourceError {
	err := r.collectAllDefinitions(items)
	if err != nil {
		return err;
	}

	err = r.bindAllNames(items)
	if err != nil {
		return err;
	}
	
	return nil;
}

func (r *resolver) collectAllDefinitions(items []parser.Item) *source.SourceError {
	for i := range items {
		item := &items[i]

		var sym parser.Symbol

		switch v := item.Data.(type) {
		case *parser.IEnum:
			sym = v.Name.Symbol
		case *parser.IRecord:
			sym = v.Name.Symbol
		default:
			continue;
		}

		if sym.IsPrimitive() {
			return &source.SourceError{
				Pos: item.Pos,
				Message: fmt.Sprintf("type '%v' shadows primitive of same name", sym),
			}
		}

		_, hasUDT := r.types[sym]
		if hasUDT {
			return &source.SourceError{
				Pos: item.Pos,
				Message: fmt.Sprintf("redifinition of type '%v'", sym),
			}
		}

		r.types[sym] = item
	}

	return nil
}

func (r *resolver) bindAllNames(items []parser.Item) *source.SourceError {
	for i := range items {
		item := &items[i]

		switch v := item.Data.(type) {
		case *parser.IRecord:
			if err := r.visitRecord(v); err != nil {
				return err;
			}
		case *parser.IService:
			if err := r.visitService(v); err != nil {
				return err;
			}
		}

	}

	return nil
}

func (r *resolver) visitRecord(record *parser.IRecord) *source.SourceError {
	for idx := range record.Fields {
		field := &record.Fields[idx]
		if err := r.visitType(&field.Type); err != nil {
			return err;
		}
	}
	return nil
}

func (r *resolver) visitService(service *parser.IService) *source.SourceError {
	for idx := range service.Operations {
		op := &service.Operations[idx];
		if err := r.visitOperation(op); err != nil {
			return err;
		}

	}
	return nil
}

func (r *resolver) visitOperation(op *parser.Operation) *source.SourceError {
	for idx := range op.Inputs {
		input := &op.Inputs[idx]
		if err := r.visitType(&input.Type); err != nil {
			return err;
		}
	}

	if err := r.visitType(&op.ResultType); err != nil {
		return err;
	}	

	return nil
}

func (r *resolver) visitType(ty *parser.Type) *source.SourceError {
	switch v := ty.Data.(type) {
	case *parser.TRef:
		return r.visitTRef(v);
	case *parser.TArray:
		return r.visitTArray(v);
	case *parser.TUnion:
		return r.visitTUnion(v);
	case *parser.TOptional:
		return r.visitTOptional(v);
	default:
		panic("unknown type variant")
	}
}

func (r *resolver) visitTArray(ty *parser.TArray) *source.SourceError {
	return r.visitType(&ty.Base)
}

func (r *resolver) visitTUnion(ty *parser.TUnion) *source.SourceError {
	for idx := range ty.Variants {
		v := &ty.Variants[idx]
		if err := r.visitType(v); err != nil {
			return err;
		}
	}
	return nil;
}

func (r *resolver) visitTOptional(ty *parser.TOptional) *source.SourceError {
	return r.visitType(&ty.Base)
}

func (r *resolver) visitTRef(ty *parser.TRef) *source.SourceError {
	sym := ty.Name.Symbol

	if sym.IsPrimitive() {
		ty.Res = sym;
		return nil;
	}

	declNode, ok := r.types[sym]
	if ok {
		ty.Res = declNode;
		return nil;
	}

	return &source.SourceError{
		Pos: ty.Name.Pos,
		Message: fmt.Sprintf("type '%v' cannot be resolved", sym),
	}
}

