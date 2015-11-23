package tftp

const (
	_            int = iota
	ModeNetascii     // RFC 1350 The TFTP Protocol (Revision 2)
	ModeOctet        // RFC 1350 The TFTP Protocol (Revision 2)
	ModeMail         // RFC 1350 The TFTP Protocol (Revision 2)
)

// Modes is a map mapping lower case mode stings to
var Modes = map[string]int{
	"netascii": ModeNetascii,
	"octet":    ModeOctet,
	"mail":     ModeMail,
}
