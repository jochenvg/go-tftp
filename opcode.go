package tftp

const (
	_                          uint16 = iota
	OpcodeReadRequest                 // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeWriteRequest                // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeData                        // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeAcknowledgment              // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeError                       // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeOptionAcknowledgment        // RFC 2347 TFTP Option Extension
)

var Opcodes = map[uint16]string{
	OpcodeReadRequest:          "RRQ",  // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeWriteRequest:         "WRQ",  // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeData:                 "DATA", // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeAcknowledgment:       "ACK",  // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeError:                "ERR",  // RFC 1350 The TFTP Protocol (Revision 2)
	OpcodeOptionAcknowledgment: "OACK", // RFC 2347 TFTP Option Extension
}
