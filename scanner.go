package tftp

import (
	"bufio"
	"bytes"
	"io"
)

// splitCString is a splitter function implementing bufio.SplitFunc.
// It splits input into null delimted strings,
func splitCString(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0); i >= 0 {
		// We have a zero terminated string.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-null terminated string. Return an error.
	if atEOF {
		return len(data), nil, ErrInvalidPacket
	}
	// Request more data.
	return 0, nil, nil
}

// NewCStringScanner returns a new bufio.Scanner which splits in tokens based
// with a null byte as delimiter.
func NewCStringScanner(reader io.Reader) (s *bufio.Scanner) {
	s = bufio.NewScanner(reader)
	s.Split(splitCString)
	return
}
