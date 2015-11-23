package tftp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"
)

var validPacketStrings = []string{
	"\x00\x01test\x00mail\x00",
	"\x00\x02test\x00netascii\x00",
	"\x00\x02test\x00octet\x00blksize\x001024\x00tsize\x000\x00timeout\x0010\x00multicast\x00\x00windowsize\x0016\x00",
	"\x00\x03\xbb\xaadata",
	"\x00\x04\xbb\xaa",
	"\x00\x05\xee\xccerror message\x00",
	"\x00\x06blksize\x001024\x00tsize\x000\x00timeout\x0010\x00multicast\x00\x00windowsize\x0016\x00",
}

var validOpcodes = []uint16{
	OpcodeReadRequest,
	OpcodeWriteRequest,
	OpcodeWriteRequest,
	OpcodeData,
	OpcodeAcknowledgment,
	OpcodeError,
	OpcodeOptionAcknowledgment,
}

var validPackets = []Packet{
	&ReadRequest{Request{
		Filename: "test",
		Mode:     "mail",
		Options:  Options{},
	}},
	&WriteRequest{Request{
		Filename: "test",
		Mode:     "netascii",
		Options:  Options{},
	}},
	&WriteRequest{Request{
		Filename: "test",
		Mode:     "octet",
		Options: Options{
			HasBlockSize:    true,
			BlockSize:       1024,
			HasTransferSize: true,
			TransferSize:    0,
			HasTimeout:      true,
			Timeout:         10,
			Multicast:       true,
			HasWindowSize:   true,
			WindowSize:      16,
		},
	}},
	&Data{
		Block: 0xbbaa,
		Data:  []byte{'d', 'a', 't', 'a'},
	},
	&Acknowledgement{
		Block: 0xbbaa,
	},
	&Error{
		ErrorCode:    0xeecc,
		ErrorMessage: "error message",
	},
	&OptionAcknowledgement{
		Options: Options{
			HasBlockSize:    true,
			BlockSize:       1024,
			HasTransferSize: true,
			TransferSize:    0,
			HasTimeout:      true,
			Timeout:         10,
			Multicast:       true,
			HasWindowSize:   true,
			WindowSize:      16,
		},
	},
}

func checkPacket(t *testing.T, p1, p2 Packet) {

	js1, _ := json.Marshal(p1)
	js2, _ := json.Marshal(p2)
	if string(js1) != string(js2) {
		fmt.Println("Packets did not match:")
		fmt.Println(string(js1))
		fmt.Println(string(js2))
		t.Error("Packets did not match")
	}

}

func TestNewPacketReadFrom_valid(t *testing.T) {
	for i, s := range validPacketStrings {
		reader := bytes.NewReader([]byte(s))
		p, n, err := NewPacketReadFrom(reader)
		// fmt.Println("Read", n, "bytes. Error is", err)
		if err != nil {
			// fmt.Println("error:", err)
			t.Error("NewPacketReaderFrom returned error")
		}
		if n != int64(len(validPacketStrings[i])) {
			t.Error("NewPackerReaderFrom did not read all bytes")
		}
		if p.Opcode() != validOpcodes[i] {
			t.Error("Packet has wrong opcode")
		}
		checkPacket(t, p, validPackets[i])
		pipereader, pipewriter := io.Pipe()
		var wg sync.WaitGroup
		wg.Add(1)
		go func(writer io.WriteCloser, p Packet, i int) {
			n, err := p.WriteTo(writer)
			if n != int64(len(validPacketStrings[i])) {
				t.Error("Packet write did not result in expected length")
			}
			if err != nil {
				t.Error("WriteTo failed with error")
			}
			writer.Close()
			wg.Done()
		}(pipewriter, p, i)
		// buf, err := ioutil.ReadAll(pipereader)
		// fmt.Printf("% x\n", buf)
		p2, n2, err2 := NewPacketReadFrom(pipereader)
		if err2 != nil {
			// fmt.Println("error:", err)
			t.Error("NewPacketReaderFrom returned error")
		}
		if n2 != int64(len(validPacketStrings[i])) {
			t.Error("NewPackerReaderFrom did not read all bytes")
		}
		if p2.Opcode() != validOpcodes[i] {
			t.Error("Packet has wrong opcode")
		}
		checkPacket(t, p, p2)
	}
}

var invalidPacketStrings = []string{
	"\x00\x07",
	"\x00\x02test",
	"\x00\x02test\x00",
	"\x00\x01test\x00mail",
	"\x00\x02test\x00octet\x00blksize\x00",
	"\x00\x02test\x00octet\x00blksize\x00WRONG\x00",
	"\x00\x02test\x00octet\x00tsize\x00WRONG\x00",
	"\x00\x02test\x00octet\x00timeout\x00WRONG\x00",
	"\x00\x02test\x00octet\x00multicast\x00WRONG\x00",
	"\x00\x02test\x00octet\x00blksize\x001024\x00tsize\x000\x00timeout\x0010\x00multicast\x00WRONG\x00windowsize\x0016\x00",
	"\x00\x02test\x00octet\x00blksize\x00WRONG\x00tsize\x000\x00timeout\x0010\x00multicast\x00\x00windowsize\x0016\x00",
	"\x00\x02test\x00octet\x00blksize\x001024\x00tsize\x000\x00timeout\x0010\x00multicast\x00\x00windowsize\x00WRONG\x00",
}

func TestNewPacketReadFrom_invalid(t *testing.T) {
	for _, s := range invalidPacketStrings {
		reader := bytes.NewReader([]byte(s))
		_, _, err := NewPacketReadFrom(reader)
		// fmt.Println("Read", n, "bytes. Error is", err)
		if err != ErrInvalidPacket {
			t.Error("NewPacketReaderFrom should return ErrInvalidPacket")
		}
	}
}
