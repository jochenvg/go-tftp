package tftp

import "errors"

// ErrInvalidPacket indicates an invalid packet was encountered
var ErrInvalidPacket = errors.New("tftp: Invalid Packet")
