package tftp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Error is an ERR TFTP Packet
type Error struct {
	ErrorCode    uint16
	ErrorMessage string
}

// Opcode gets the opcode
func (e *Error) Opcode() uint16 {
	return OpcodeError
}

// ReadFrom implements the io.ReaderFrom interface
func (e *Error) ReadFrom(reader io.Reader) (n int64, err error) {
	ecr := NewErrorCountReader(reader)
	binary.Read(ecr, binary.BigEndian, &e.ErrorCode)
	scanner := NewCStringScanner(ecr)
	if scanner.Scan() {
		e.ErrorMessage = scanner.Text()
	}
	err = scanner.Err()
	if err == nil && ecr.Err() != io.EOF {
		err = ecr.Err()
	}
	return ecr.Count(), err
}

// WriteTo implements the io.WriterTo interface
func (e *Error) WriteTo(writer io.Writer) (n int64, err error) {
	ecw := NewErrorCountWriter(writer)
	binary.Write(ecw, binary.BigEndian, OpcodeError)
	binary.Write(ecw, binary.BigEndian, e.ErrorCode)
	fmt.Fprintf(ecw, "%s\x00", e.ErrorMessage)
	return ecw.Count(), ecw.Err()
}
