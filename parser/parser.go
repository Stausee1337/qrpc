package parser

import (
	"fmt"

	"github.com/stausee1337/qrpc/lexer"
	"github.com/stausee1337/qrpc/source"
)

type parser struct {
	stream 	 []lexer.Token
	position uint
}

func ParseTokenStream(stream []lexer.Token) ([]Item, *source.SourceError) {
	p := parser { stream: stream }

	items := make([]Item, 0)
	for !p.isEOF() {
		item, err := p.parseItem();
		if err != nil { return nil, err }
		items = append(items, item)
	}
	return items, nil
}

func (p *parser) parseItem() (Item, *source.SourceError) {
	tok := p.current()

	switch tok.Kind {
	case lexer.KeywordEnum:
		return bindItem(p, p.parseEnum);
	case lexer.KeywordRecord:
		return bindItem(p, p.parseRecord);
	case lexer.KeywordService:
		return bindItem(p, p.parseService);
	default:
		return Item{}, p.unexpected()
	}
}

func (p *parser) parseEnum() (*IEnum, *source.SourceError) {
	name, err := p.expect(lexer.Ident)
	if err != nil { return nil, err }

	_, err = p.expect(lexer.LBrace)
	if err != nil { return nil, err }

	variants := make([]Ident, 0)
	for p.current().Kind != lexer.RBrace && !p.isEOF() {
		variant, err := p.expect(lexer.Ident)
		if err != nil { return nil, err }

		tok := p.current()
		if tok.Kind != lexer.RBrace && tok.Kind != lexer.Comma {
			return nil, p.unexpected()
		}
		p.skipIf(lexer.Comma);

		variants = append(variants, makeIdent(variant))
	}

	_, err = p.expect(lexer.RBrace)
	if err != nil { return nil, err }

	return &IEnum{
		Name: makeIdent(name),
		Variants: variants,
	}, nil
}

func (p *parser) parseRecord() (*IRecord, *source.SourceError) {
	name, err := p.expect(lexer.Ident)
	if err != nil { return nil, err }

	_, err = p.expect(lexer.LBrace)
	if err != nil { return nil, err }

	fields := make([]Field, 0)
	for p.current().Kind != lexer.RBrace && !p.isEOF() {
		field, err := p.expect(lexer.Ident)
		if err != nil { return nil, err }

		_, err = p.expect(lexer.Colon);
		if err != nil { return nil, err }

		ty, err := p.parseType()
		if err != nil { return nil, err }

		_, err = p.expect(lexer.Semicolon);
		if err != nil { return nil, err }

		fields = append(fields, Field{
			Name: makeIdent(field),
			Type: ty,
		})
	}

	_, err = p.expect(lexer.RBrace)
	if err != nil { return nil, err }

	return &IRecord{
		Name: makeIdent(name),
		Fields: fields,
	}, nil
}

func (p *parser) parseService() (*IService, *source.SourceError) {
	name, err := p.expect(lexer.Ident)
	if err != nil { return nil, err }

	_, err = p.expect(lexer.LBrace)
	if err != nil { return nil, err }

	operations := make([]Operation, 0)
	for p.current().Kind != lexer.RBrace && !p.isEOF() {
		operation, err := p.parseOperation()
		if err != nil { return nil, err }

		_, err = p.expect(lexer.Semicolon);
		if err != nil { return nil, err }

		operations = append(operations, operation)
	}

	_, err = p.expect(lexer.RBrace)
	if err != nil { return nil, err }

	return &IService{
		Name: makeIdent(name),
		Operations: operations,
	}, nil
}

func bindItem[K I](p *parser, body func() (K, *source.SourceError)) (Item, *source.SourceError) {
	start := p.advance()

	data, err := body()
	end := &p.stream[p.position - 1]

	if err != nil {
		return Item{}, err
	}

	return Item{
		Pos: source.MorphPosition(&start.Pos, &end.Pos),
		Data: data,
	}, nil
}

func (p *parser) parseOperation() (Operation, *source.SourceError) {
	start := p.current()
	var kind OperationKind

	switch start.Kind {
	case lexer.KeywordQuery:
		kind = OperationQuery
	case lexer.KeywordMutation:
		kind = OperationMutation
	default:
		return Operation{}, p.unexpected()
	}
	p.advance()

	name, err := p.expect(lexer.Ident)
	if err != nil { return Operation{}, err }

	_, err = p.expect(lexer.LParen)
	if err != nil { return Operation{}, err }

	inputs := make([]NamedInput, 0)
	for p.current().Kind != lexer.RParen && !p.isEOF() {
		input, err := p.parseNamedInput()
		if err != nil { return Operation{}, err }

		tok := p.current()
		if tok.Kind != lexer.RParen && tok.Kind != lexer.Comma {
			return Operation{}, p.unexpected()
		}
		p.skipIf(lexer.Comma);

		inputs = append(inputs, input)
	}

	_, err = p.expect(lexer.RParen)
	if err != nil { return Operation{}, err }

	_, err = p.expect(lexer.Colon)
	if err != nil { return Operation{}, err }

	resultType, err := p.parseType()
	if err != nil { return Operation{}, err }

	return Operation{
		Pos: source.MorphPosition(&start.Pos, &resultType.Pos),
		Kind: kind,
		Name: makeIdent(name),
		Inputs: inputs,
		ResultType: resultType,
	}, nil
}

func (p *parser) parseNamedInput() (NamedInput, *source.SourceError) {
	name, err := p.expect(lexer.Ident)
	if err != nil { return NamedInput{}, err }

	_, err = p.expect(lexer.Colon)
	if err != nil { return NamedInput{}, err }

	ty, err := p.parseType()
	if err != nil { return NamedInput{}, err }

	return NamedInput{
		Pos: source.MorphPosition(&name.Pos, &ty.Pos),
		Name: makeIdent(name),
		Type: ty,
	}, nil;
}

func (p *parser) parseType() (Type, *source.SourceError) {
	base, err := p.parsePostfixType();
	if err != nil { return Type{}, err }

	if p.current().Kind != lexer.VBar {
		return base, nil;
	}

	variants := []Type{ base }
	for p.skipIf(lexer.VBar) != nil {
		variant, err := p.parsePostfixType();
		if err != nil { return Type{}, err }

		variants = append(variants, variant)
	}

	return Type{
		Pos: source.MorphPosition(&base.Pos, &variants[len(variants) - 1].Pos),
		Data: &TUnion{ Variants: variants },
	}, nil
}

func (p *parser) parsePostfixType() (Type, *source.SourceError) {
	ty, err := p.parseRefType()
	if err != nil { return Type{}, err }

	if q := p.skipIf(lexer.Question); q != nil {
		ty = Type {
			Pos: source.MorphPosition(&ty.Pos, &q.Pos),
			Data: &TOptional{ Base: ty },
		}
	}

	if p.skipIf(lexer.LBracket) != nil {
		tok, err := p.expect(lexer.RBracket)
		if err != nil { return Type{}, err }

		ty = Type {
			Pos: source.MorphPosition(&ty.Pos, &tok.Pos),
			Data: &TArray{ Base: ty },
		}
	}

	return ty, nil
}

func (p *parser) parseRefType() (Type, *source.SourceError) {
	ident, err := p.expect(lexer.Ident)
	if err != nil { return Type{}, err }

	name := makeIdent(ident)
	return Type {
		Pos: name.Pos,
		Data: &TRef{ Name: name },
	}, nil
}

func (p *parser) skipIf(kind lexer.Kind) *lexer.Token {
	tok := p.current()
	if tok.Kind == kind {
		p.advance()
		return tok
	}
	return nil
}

func (p *parser) expect(kind lexer.Kind) (*lexer.Token, *source.SourceError) {
	tok := p.current();

	if tok.Kind == kind {
		p.advance();
		return tok, nil;
	}

	return nil, &source.SourceError{
		Pos: tok.Pos,
		Message: fmt.Sprintf("expected '%v', found '%v'", kind, tok.Kind),
	};
}

func (p *parser) unexpected() *source.SourceError {
	tok := p.current();

	return &source.SourceError{
		Pos: tok.Pos,
		Message: fmt.Sprintf("unexpected '%v'", tok.Kind),
	}
}

func (p *parser) advance() *lexer.Token {
	current := &p.stream[p.position]
	if current.Kind != lexer.EOF {
		p.position++;
	}
	return current
}

func (p *parser) current() *lexer.Token {
	return &p.stream[p.position]
}

func (p *parser) isEOF() bool {
	return p.current().Kind == lexer.EOF
}

