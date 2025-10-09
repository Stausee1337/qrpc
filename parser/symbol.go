package parser

type Symbol uint

const (
	SYMstring Symbol = iota
	SYMbool
	SYMempty
	SYMfloat
	SYMint
	SYMuint
	SYMuuid
)

var symbolCache = map[string]Symbol {
	"string": SYMstring,
	"bool": SYMbool,
	"empty": SYMempty,
	"float": SYMfloat,
	"int": SYMint,
	"uint": SYMuint,
	"UUID": SYMuuid,
}

func GetSymbol(str string) Symbol {
	sym, ok := symbolCache[str]
	if ok {
		return sym
	}

	sym = Symbol(len(symbolCache))
	symbolCache[str] = sym
	return sym
}


