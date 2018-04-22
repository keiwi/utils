package file

// stats tracks the number of output lines and bytes written.
type stats struct {
	lines int64
	bytes int64
}

// Lines returns the number of lines written.
func (s *stats) Lines() int64 {
	return s.lines
}

// Bytes returns the number of bytes written.
func (s *stats) Bytes() int64 {
	return s.bytes
}
