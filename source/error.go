package source

type SyntaxError struct {
	Pos		Position
	Message string
}

func (s *SyntaxError) Error() string {
	return s.Message;
}

