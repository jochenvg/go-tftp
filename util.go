package tftp

import "io"

func sliceTo0(p []byte) ([]byte, []byte) {
	for i, b := range p {
		if b == 0 {
			return p[0:i], p[i+1:]
		}
	}
	return p[:], nil
}

// ErrorCountWriter embeds an io.Writer and counts the total number of bytes written to it over its lifetime.
// It also tracks the error state, which is defined as the first write error encountered.
type ErrorCountWriter struct {
	io.Writer
	nn  int64
	err error
}

// NewErrorCountWriter returns a new ErrorCountWriter
func NewErrorCountWriter(writer io.Writer) *ErrorCountWriter {
	return &ErrorCountWriter{
		Writer: writer,
	}
}

// Write writes the contents of p if no previous write errors occurred.
// It returns the number of bytes written in the current write operation
// as well as the error condition on the writer
func (e *ErrorCountWriter) Write(p []byte) (n int, err error) {
	// If no write errors encountered, Write
	if e.err == nil {
		n, err = e.Writer.Write(p)
		e.nn += int64(n)
		e.err = err
	}
	return n, e.err
}

// Count returns the cumulated number of bytes written
func (e *ErrorCountWriter) Count() int64 {
	return e.nn
}

// Err returns the error state
func (e *ErrorCountWriter) Err() error {
	return e.err
}

// Reset resets the count and error state
func (e *ErrorCountWriter) Reset() {
	e.err = nil
	e.nn = 0
}

// ErrorCountReader embeds a io.Reader and counts the total number of bytes written to it over its lifetime
type ErrorCountReader struct {
	io.Reader
	nn  int64
	err error
}

// NewErrorCountReader returns a new ErrorCountReader
func NewErrorCountReader(reader io.Reader) *ErrorCountReader {
	return &ErrorCountReader{
		Reader: reader,
	}
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained.  The return value n is the number of bytes read.  If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (e *ErrorCountReader) Read(p []byte) (n int, err error) {
	if e.err == nil {
		n, err = e.Reader.Read(p)
		e.nn += int64(n)
		e.err = err
	}
	return n, e.err
}

// Count returns the cumulated number of bytes read
func (e *ErrorCountReader) Count() int64 {
	return e.nn
}

// Err returns the error state
func (e *ErrorCountReader) Err() error {
	return e.err
}

// Reset resets the count and error state
func (e *ErrorCountReader) Reset() {
	e.err = nil
	e.nn = 0
}
