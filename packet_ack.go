package tftp

import (
	"encoding/binary"
	"io"
)

// Acknowledgement is a DATA TFTP Packet
type Acknowledgement struct {
	Block uint16
}

// Opcode gets the opcode
func (a *Acknowledgement) Opcode() uint16 {
	return OpcodeAcknowledgment
}

// ReadFrom implements the io.ReaderFrom interface
func (a *Acknowledgement) ReadFrom(reader io.Reader) (n int64, err error) {
	ecr := NewErrorCountReader(reader)
	err = binary.Read(ecr, binary.BigEndian, &a.Block)
	return ecr.Count(), ecr.Err()
}

// WriteTo implements the io.WriterTo interface
func (a *Acknowledgement) WriteTo(writer io.Writer) (n int64, err error) {
	ecw := NewErrorCountWriter(writer)
	binary.Write(ecw, binary.BigEndian, OpcodeAcknowledgment)
	binary.Write(ecw, binary.BigEndian, a.Block)
	return ecw.Count(), ecw.Err()
}
