package tftp

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

// Data is a DATA TFTP Packet
type Data struct {
	Block uint16
	Data  []byte
}

// Opcode gets the opcode
func (d *Data) Opcode() uint16 {
	return OpcodeData
}

// ReadFrom implements the io.ReaderFrom interface
func (d *Data) ReadFrom(reader io.Reader) (n int64, err error) {
	ecr := NewErrorCountReader(reader)
	binary.Read(ecr, binary.BigEndian, &d.Block)
	d.Data, err = ioutil.ReadAll(ecr)
	if err == nil && ecr.Err() != io.EOF {
		err = ecr.Err()
	}
	return ecr.Count(), err
}

// WriteTo writes a DATA packet in byte form
func (d *Data) WriteTo(writer io.Writer) (n int64, err error) {
	ecw := NewErrorCountWriter(writer)
	binary.Write(ecw, binary.BigEndian, OpcodeData)
	binary.Write(ecw, binary.BigEndian, d.Block)
	ecw.Write(d.Data)
	return ecw.Count(), ecw.Err()
}
