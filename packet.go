package tftp

import (
	"encoding/binary"
	"io"
)

// Packet is a TFTP packet
type Packet interface {
	// Opcode gets the packet opcode
	Opcode() uint16
	// ReadFrom implements the io.ReaderFrom interface
	ReadFrom(reader io.Reader) (n int64, err error)
	// WriteTo implements the io.WriterTo interface
	WriteTo(writer io.Writer) (n int64, err error)
}

// NewPacketReadFrom returns a new Packet by reading it from reader
func NewPacketReadFrom(reader io.Reader) (p Packet, n int64, err error) {
	var opcode uint16
	var packet Packet
	if err = binary.Read(reader, binary.BigEndian, &opcode); err != nil {
		return
	}
	switch opcode {
	case OpcodeReadRequest:
		packet = &ReadRequest{}
	case OpcodeWriteRequest:
		packet = &WriteRequest{}
	case OpcodeData:
		packet = &Data{}
	case OpcodeAcknowledgment:
		packet = &Acknowledgement{}
	case OpcodeError:
		packet = &Error{}
	case OpcodeOptionAcknowledgment:
		packet = &OptionAcknowledgement{}
	default:
		err = ErrInvalidPacket
		return
	}
	p = packet
	n, err = packet.ReadFrom(reader)
	n += int64(binary.Size(opcode))
	return
}
