package tftp

import "net"

func ExampleNewServer() {
	server := Server{}
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1")
	server.Listen(addr)
	defer server.Close()
	server.Serve()
}

// func TestNewServer(t *testing.T) {
// 	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:6969")
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	server := Server{}
// 	err = server.Listen(addr)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}
// 	defer server.Close()
// 	server.Serve()
// }
