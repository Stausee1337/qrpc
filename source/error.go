package source

type SourceError struct {
	Pos		Position
	Message string
}

func (s *SourceError) Error() string {
	return s.Message;
}

