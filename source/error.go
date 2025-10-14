package source

import (
	"fmt"
	"os"
)

type SourceError struct {
	Pos     Position
	Message string
}

func (s *SourceError) Error() string {
	return s.Message
}

func RenderSourceError(err *SourceError) {
	fmt.Fprintf(
		os.Stderr,
		"ERROR: %v:%v:%v: %v\n",
		err.Pos.File, err.Pos.Lineno, err.Pos.Column,
		err.Message,
	)
}
