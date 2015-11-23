package tftp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Options represents TFTP options as per:
// RFC 2347 TFTP Option Extension
// RFC 2348 TFTP Blocksize Option,
// RFC 2349 TFTP Timeout Interval and Transfer Size Options,
// RFC 2090 TFTP Multicast Option,
// RFC 7440 TFTP Windowsize Option
type Options struct {
	HasBlockSize    bool
	BlockSize       int // RFC 2348 TFTP Blocksize Option
	HasTimeout      bool
	Timeout         int // RFC 2349 TFTP Timeout Interval and Transfer Size Options
	HasTransferSize bool
	TransferSize    int  // RFC 2349 TFTP Timeout Interval and Transfer Size Options
	Multicast       bool // RFC 2090 TFTP Multicast Option
	HasWindowSize   bool
	WindowSize      int // RFC 7440 TFTP Windowsize Option
}

// ScanFrom scans from a bufio.Scanner
func (o *Options) ScanFrom(scanner *bufio.Scanner) (err error) {
	// Loop and get the first token
	for scanner.Scan() {
		var name, value string
		var e error

		name = strings.ToLower(scanner.Text())
		if !scanner.Scan() {
			err = ErrInvalidPacket
			break
		}
		value = scanner.Text()
		// RFC 2347 TFTP Option Extension specifies null terminated strings
		switch name {
		case "blksize": // RFC 2348
			o.HasBlockSize = true
			o.BlockSize, e = strconv.Atoi(string(value))
			if e != nil || o.BlockSize < 8 || o.BlockSize > 65464 {
				err = ErrInvalidPacket
			}
		case "timeout": // RFC 2349
			o.HasTimeout = true
			o.Timeout, e = strconv.Atoi(string(value))
			if e != nil {
				err = ErrInvalidPacket
			}
		case "tsize": // RFC 2349
			o.HasTransferSize = true
			o.TransferSize, e = strconv.Atoi(string(value))
			if e != nil || o.TransferSize < 0 {
				err = ErrInvalidPacket
			}
		case "multicast": // RFC 2090
			if len(value) == 0 {
				o.Multicast = true
			} else {
				err = ErrInvalidPacket
			}
		case "windowsize": // RFC 7440
			o.HasWindowSize = true
			o.WindowSize, e = strconv.Atoi(string(value))
			if err != nil || o.WindowSize < 1 || o.WindowSize > 65535 {
				err = ErrInvalidPacket
			}
		}
	}
	if err == nil {
		err = scanner.Err()
	}
	return
}

// WriteTo implements the io.WriterTo interface
func (o *Options) WriteTo(writer io.Writer) (n int64, err error) {
	ecw := NewErrorCountWriter(writer)
	if o.HasBlockSize {
		fmt.Fprintf(ecw, "blksize\x00%d\x00", o.BlockSize)
	}
	if o.HasTimeout {
		fmt.Fprintf(ecw, "timeout\x00%d\x00", o.Timeout)
	}
	if o.HasTransferSize {
		fmt.Fprintf(ecw, "tsize\x00%d\x00", o.TransferSize)
	}
	if o.Multicast {
		fmt.Fprintf(ecw, "multicast\x00\x00")
	}
	if o.HasWindowSize {
		fmt.Fprintf(ecw, "windowsize\x00%d\x00", o.WindowSize)
	}
	return ecw.Count(), ecw.Err()
}
