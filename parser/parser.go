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

func ParseTokenStream(stream []lexer.Token) *source.SyntaxError {
	p := parser { stream: stream }

	for !p.isEOF() {
		_, err := p.parseItem();
		if err != nil { return err }
	}
	return nil
}

func (p *parser) parseItem() (Item, *source.SyntaxError) {
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

func (p *parser) parseEnum() (*IEnum, *source.SyntaxError) {
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

func (p *parser) parseRecord() (*IRecord, *source.SyntaxError) {
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

func (p *parser) parseService() (*IService, *source.SyntaxError) {
	name, err := p.expect(lexer.Ident)
	if err != nil { return nil, err }

	_, err = p.expect(lexer.LBrace)
	if err != nil { return nil, err }

	fields := make([]Operation, 0)
	for p.current().Kind != lexer.RBrace && !p.isEOF() {
		operation, err := p.parseOperation()
		if err != nil { return nil, err }

		_, err = p.expect(lexer.Semicolon);
		if err != nil { return nil, err }

		fields = append(fields, operation)
	}

	_, err = p.expect(lexer.RBrace)
	if err != nil { return nil, err }

	return &IService{
		Name: makeIdent(name),
	}, nil
}

func bindItem[K I](p *parser, body func() (K, *source.SyntaxError)) (Item, *source.SyntaxError) {
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

func (p *parser) parseOperation() (Operation, *source.SyntaxError) {
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

	end, err := p.expect(lexer.RParen)
	if err != nil { return Operation{}, err }

	endPos := &end.Pos
	var resultType *Type

	if p.skipIf(lexer.Colon) != nil {
		r, err := p.parseType()
		if err != nil { return Operation{}, err }
		resultType = &r
		endPos = &resultType.Pos
	}

	return Operation{
		Pos: source.MorphPosition(&start.Pos, endPos),
		Kind: kind,
		Name: makeIdent(name),
		Inputs: inputs,
		ResultType: resultType,
	}, nil
}

func (p *parser) parseNamedInput() (NamedInput, *source.SyntaxError) {
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

func (p *parser) parseType() (Type, *source.SyntaxError) {
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

func (p *parser) parsePostfixType() (Type, *source.SyntaxError) {
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

func (p *parser) parseRefType() (Type, *source.SyntaxError) {
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

func (p *parser) expect(kind lexer.Kind) (*lexer.Token, *source.SyntaxError) {
	tok := p.current();

	if tok.Kind == kind {
		p.advance();
		return tok, nil;
	}

	return nil, &source.SyntaxError{
		Pos: tok.Pos,
		Message: fmt.Sprintf("expected %q, found %q", kind, tok.Kind),
	};
}

func (p *parser) unexpected() *source.SyntaxError {
	tok := p.current();

	return &source.SyntaxError{
		Pos: tok.Pos,
		Message: fmt.Sprintf("unexpected %q", tok.Kind),
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

