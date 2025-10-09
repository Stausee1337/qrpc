package source

type Position struct {
	Start  uint
	End    uint
	Lineno uint
	Column uint
	File   string
}

func MorphPosition(start *Position, end *Position) Position {
	return Position{
		Start: start.Start,
		End: end.End,
		Lineno: start.Lineno,
		Column: start.Column,
		File: start.File,
	}
}

