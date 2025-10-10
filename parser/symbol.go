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
	numPrimitiveSym
)

type indexSet struct {
	hashmap map[string]Symbol
	list	[]string
}

func createIndexSet(initialSymbols ...string) indexSet {
	set := indexSet {
		hashmap: map[string]Symbol{},
		list: []string{},
	}

	for _, string := range initialSymbols {
		set.insert(string);
	}

	return set;
}

func (i *indexSet) insert(str string) Symbol {
	sym := Symbol(len(i.list))
	i.hashmap[str] = sym
	i.list = append(i.list, str)
	return sym
}

func (i *indexSet) lookup(str string) (Symbol, bool) {
	sym, ok := i.hashmap[str]
	return sym, ok
}

func (i *indexSet) getString(sym Symbol) string {
	return i.list[sym]
}

var symbolCache indexSet

func init() {
	symbolCache = createIndexSet(
		"string",
		"bool",
		"empty",
		"float",
		"int",
		"uint",
	)
}

func InternSymbol(str string) Symbol {
	sym, ok := symbolCache.lookup(str)
	if ok {
		return sym
	}

	return symbolCache.insert(str)
}

func (sym Symbol) String() string {
	return symbolCache.getString(sym)
}

func (sym Symbol) IsPrimitive() bool {
	return sym < numPrimitiveSym
}

