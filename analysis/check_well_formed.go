package analysis

import (
	"fmt"

	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

func runWellFormedChecks(items []parser.Item) *source.SourceError {
	for _, item := range items {
		switch v := item.Data.(type) {
		case *parser.IEnum:
			if err := checkEnumWellFormed(v, &item.Pos); err != nil {
				return err;
			}
		case *parser.IRecord:
			if err := checkRecordWellFormed(v, &item.Pos); err != nil {
				return err;
			}
		case *parser.IService:
			if err := checkServiceWellFormed(v, &item.Pos); err != nil {
				return err;
			}
		}
	}

	return nil;
}

func checkEnumWellFormed(enum *parser.IEnum, pos *source.Position) *source.SourceError {
	if len(enum.Variants) == 0 {
		return &source.SourceError{
			Pos: *pos,
			Message: fmt.Sprintf("enum '%v' must not be empty", enum.Name.Symbol),
		}
	}

 	variantMap := map[parser.Symbol]bool{}
	for _, variant := range enum.Variants {
		_, hasVariant := variantMap[variant.Symbol]
		if hasVariant {
			return &source.SourceError{
				Pos: variant.Pos,
				Message: fmt.Sprintf("variant '%v' appears multiple times in enum '%v'", variant.Symbol, enum.Name.Symbol),
			}
		}
		variantMap[variant.Symbol] = true
	}

	return nil;
}

func checkRecordWellFormed(record *parser.IRecord, pos *source.Position) (*source.SourceError) {
	if len(record.Fields) == 0 {
		return &source.SourceError{
			Pos: *pos,
			Message: fmt.Sprintf("record '%v' must not be empty", record.Name.Symbol),
		}
	}

 	fieldMap := map[parser.Symbol]bool{}
	for _, field := range record.Fields {
		_, hasField := fieldMap[field.Name.Symbol]
		if hasField {
			return &source.SourceError{
				Pos: field.Name.Pos,
				Message: fmt.Sprintf("field '%v' appears multiple times in record '%v'", field.Name.Symbol, record.Name.Symbol),
			}
		}
		fieldMap[field.Name.Symbol] = true
	}

	return nil;
}

func checkServiceWellFormed(service *parser.IService, pos *source.Position) (*source.SourceError) {
	if len(service.Operations) == 0 {
		return &source.SourceError{
			Pos: *pos,
			Message: fmt.Sprintf("service '%v' must not be empty", service.Name.Symbol),
		}
	}

 	opMap := map[parser.Symbol]bool{}
	for idx := range service.Operations {
		op := &service.Operations[idx]
		_, hasOp := opMap[op.Name.Symbol]
		if hasOp {
			return &source.SourceError{
				Pos: op.Name.Pos,
				Message: fmt.Sprintf("operation '%v' appears multiple times in service '%v'", op.Name.Symbol, service.Name.Symbol),
			}
		}
		opMap[op.Name.Symbol] = true

		if err := checkOperationWellFormed(op, service.Name.Symbol); err != nil {
			return err;
		}
	}

	return nil;

}

func checkOperationWellFormed(op *parser.Operation, serviceSymbol parser.Symbol) *source.SourceError {

 	inputMap := map[parser.Symbol]bool{}
	for idx := range op.Inputs {
		input := &op.Inputs[idx]
		_, hasInput := inputMap[input.Name.Symbol]
		if hasInput {
			return &source.SourceError{
				Pos: input.Name.Pos,
				Message: fmt.Sprintf(
					"input '%v' appears multiple times in operation '%v.%v'",
					input.Name.Symbol,
					serviceSymbol,
					op.Name.Symbol,
				),
			}
		}
		inputMap[input.Name.Symbol] = true
	}

	if op.Kind != parser.OperationMutation {
		return nil;
	}

	if len(op.Inputs) == 0 {
		return &source.SourceError{
			Pos: op.Name.Pos,
			Message: fmt.Sprintf("mutation '%v' in service '%v' must have at least one input", op.Name.Symbol, serviceSymbol),
		}
	}

	return nil;
}

