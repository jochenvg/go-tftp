package tftp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Request is a TFTP RRQ Packet
type Request struct {
	Filename string
	Mode     string
	Options  Options
}

// ReadFrom implements the io.ReaderFrom interface
func (rq *Request) ReadFrom(reader io.Reader) (n int64, err error) {
	// RFC 1350 The TFTP Protocol (Revision 2)
	ecr := NewErrorCountReader(reader)
	scanner := NewCStringScanner(ecr)

	if scanner.Scan() {
		rq.Filename = scanner.Text()
	}
	if scanner.Scan() {
		rq.Mode = scanner.Text()
	}
	err = rq.Options.ScanFrom(scanner)
	if err == nil {
		err = scanner.Err()
	}
	if err == nil && ecr.Err() != io.EOF {
		err = ecr.Err()
	}
	if rq.Filename == "" || rq.Mode == "" {
		err = ErrInvalidPacket
	}
	return ecr.Count(), err
}

// WriteToWithOpcode implements the io.WriterFrom interface with an additional opcode parameter
func (rq *Request) opcodeWriteTo(writer io.Writer, opcode uint16) (n int64, err error) {
	ecw := NewErrorCountWriter(writer)
	binary.Write(ecw, binary.BigEndian, opcode)
	fmt.Fprintf(ecw, "%s\x00%s\x00", rq.Filename, rq.Mode)
	rq.Options.WriteTo(ecw)
	return ecw.Count(), ecw.Err()
}

// ReadRequest is a TFTP RRQ packet
type ReadRequest struct {
	Request
}

// Opcode gets the opcode
func (rq *ReadRequest) Opcode() uint16 {
	return OpcodeReadRequest
}

// WriteTo implements the io.WriterTo interface
func (rq *ReadRequest) WriteTo(writer io.Writer) (n int64, err error) {
	n, err = rq.opcodeWriteTo(writer, OpcodeReadRequest)
	return
}

// WriteRequest is a TFTP WRQ packet
type WriteRequest struct {
	Request
}

// Opcode gets the opcode
func (rq *WriteRequest) Opcode() uint16 {
	return OpcodeWriteRequest
}

// WriteTo implements the io.WriterTo interface
func (rq *WriteRequest) WriteTo(writer io.Writer) (n int64, err error) {
	n, err = rq.opcodeWriteTo(writer, OpcodeWriteRequest)
	return
}
