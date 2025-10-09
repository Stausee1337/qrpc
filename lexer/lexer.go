package lexer


type LexError struct {
	Span Span
}

func (e *LexError) Error() string {
	return "unexpected token"
}

type lexer struct {
	source string
	start  uint
	end    uint
	bol    uint
	lineno uint
}

func LexToStream(s string) ([]Token, error) {
	l := lexer{source: s, lineno: 1}

	slice := make([]Token, 0)
	for {
		tok := l.lex()
		// fmt.Printf("%v, %v\n", l.end, uint(len(l.source)));
		if tok.Kind == Invalid {
			return nil, &LexError{Span: tok.Span}
		}
		slice = append(slice, tok)
		if tok.Kind == EOF {
			break
		}
	}
	return slice, nil
}

func (l *lexer) lex() Token {
	l.skipWhitespace()

	l.start = l.end
	if l.start >= uint(len(l.source)) {
		return Token{
			Span:  l.position(),
			Kind:  EOF,
			Value: "",
		}
	}

	switch l.source[l.end] {
	case '(':
		l.end++
		return l.bindEmptyToken(LParen)
	case ')':
		l.end++
		return l.bindEmptyToken(RParen)
	case '[':
		l.end++
		return l.bindEmptyToken(LBracket)
	case ']':
		l.end++
		return l.bindEmptyToken(RBracket)
	case '{':
		l.end++
		return l.bindEmptyToken(LBrace)
	case '}':
		l.end++
		return l.bindEmptyToken(RBrace)
	case ':':
		l.end++
		return l.bindEmptyToken(Colon)
	case ',':
		l.end++
		return l.bindEmptyToken(Comma)
	case ';':
		l.end++
		return l.bindEmptyToken(Semicolon)
	case '?':
		l.end++
		return l.bindEmptyToken(Question)
	case '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return l.lexIdentOrKeyword()
	default:
		return l.bindEmptyToken(Invalid)
	}

}

var keywords = map[string]Kind{
	"enum":     KeywordEnum,
	"record":   KeywordRecord,
	"service":  KeywordService,
	"query":    KeywordQuery,
	"mutation": KeywordMutation,
}

func (l *lexer) lexIdentOrKeyword() Token {
	for {
		l.end++

		if l.end >= uint(len(l.source)) || !isNameCharacter(l.source[l.end]) {
			break
		}
	}

	value := l.source[l.start:l.end]
	kind, ok := keywords[value]

	if ok {
		return l.bindEmptyToken(kind)
	}

	return l.bindValueToken(Ident, value)
}

func isNameCharacter(r byte) bool {
	return (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_'
}

func (l *lexer) skipWhitespace() {
	for l.end < uint(len(l.source)) {
		switch l.source[l.end] {
		case '\t', ' ':
			l.end++
		case '\n':
			l.end++
			l.lineno++
			l.bol = l.end
		case '\r':
			l.end++
			l.lineno++

			if l.end < uint(len(l.source)) && l.source[l.end] == '\n' {
				l.end++
			}

			l.bol = l.end
		default:
			return
		}
	}
}

func (l *lexer) bindEmptyToken(kind Kind) Token {
	return Token{
		Span:  l.position(),
		Kind:  kind,
		Value: "",
	}
}

func (l *lexer) bindValueToken(kind Kind, value string) Token {
	return Token{
		Span:  l.position(),
		Kind:  kind,
		Value: value,
	}
}

func (l *lexer) position() Span {
	return Span{
		Start:  l.start,
		End:    l.end,
		Lineno: l.lineno,
		Column: (l.start - l.bol) + 1,
	}
}
