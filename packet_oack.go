package tftp

import (
	"encoding/binary"
	"io"
)

// OptionAcknowledgement is a TFTP Request Packet
type OptionAcknowledgement struct {
	Options Options
}

// Opcode gets the opcode
func (o *OptionAcknowledgement) Opcode() uint16 {
	return OpcodeOptionAcknowledgment
}

// ReadFrom implements the io.ReaderFrom interface
func (o *OptionAcknowledgement) ReadFrom(reader io.Reader) (n int64, err error) {
	ecr := NewErrorCountReader(reader)
	scanner := NewCStringScanner(ecr)
	// RFC 2347 TFTP Option Extension
	// Options
	err = o.Options.ScanFrom(scanner)
	if err == nil && ecr.Err() != io.EOF {
		err = ecr.Err()
	}
	return ecr.Count(), err
}

// WriteTo implements the io.WriterTo interface
func (o *OptionAcknowledgement) WriteTo(writer io.Writer) (n int64, err error) {
	ecw := NewErrorCountWriter(writer)
	binary.Write(ecw, binary.BigEndian, OpcodeOptionAcknowledgment)
	o.Options.WriteTo(ecw)
	return ecw.Count(), err
}
