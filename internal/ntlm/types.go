package ntlm

const (
	nTLM_NEGOTIATE_UNICODE uint32 = 1 << iota
	nTLM_NEGOTIATE_OEM     uint32 = 1 << iota
	nTLM_REQUEST_TARGET           = 1 << iota
	_
	nTLM_NEGOTIATE_SIGN     uint32 = 1 << iota
	nTLM_NEGOTIATE_SEAL     uint32 = 1 << iota
	nTLM_NEGOTIATE_DATAGRAM uint32 = 1 << iota
	nTLM_NEGOTIATE_LM_KEY   uint32 = 1 << iota
	_
	nTLM_NEGOTIATE_NTLM uint32 = 1 << iota
	_
	nTLM_NEGOTIATE_ANONYMOUS            uint32 = 1 << iota
	nTLM_NEGOTIATE_DOMAIN_SUPPLIED      uint32 = 1 << iota
	nTLM_NEGOTIATE_WORKSTATION_SUPPLIED uint32 = 1 << iota
	_
	nTLM_NEGOTIATE_ALWAYS_SIGN uint32 = 1 << iota
	nTLM_TARGET_TYPE_DOMAIN    uint32 = 1 << iota
	nTLM_TARGET_TYPE_SERVER    uint32 = 1 << iota
	_
	nTLM_NEGOTIATE_EXTENDED_SESSION_SECURITY uint32 = 1 << iota
	nTLM_NEGOTIATE_IDENTIFY                  uint32 = 1 << iota
	_
	nTLM_NEGOTIATE_REQUEST_NON_NT_KEY uint32 = 1 << iota
	nTLM_NEGOITATE_TARGET_INFO        uint32 = 1 << iota
	_
	nTLM_NEGOTIATE_VERSION uint32 = 1 << iota
	_
	_
	_
	nTLM_NEGOTIATE_128      uint32 = 1 << iota
	nTLM_NEGOTIATE_KEY_EXCH uint32 = 1 << iota
	nTLM_NEGOTIATE_56       uint32 = 1 << iota
)

const (
	avID_MsvAvEOL uint16 = iota
	avID_MsvAvNbComputerName
	avID_MsvAvNbDomainName
	avID_MsvAvDnsComputerName
	avID_MsvAvDnsDomainName
	avID_MsvAvDnsTreeName
	avID_MsvAvFlags
	avID_MsvAvTimestamp
	avID_MsvAvSingleHost
	avID_MsvAvTargetName
	avID_MsvChannelBindings
)

var (
	nTLM_MESSAGE_SIGNATURE [8]byte = [8]byte{'N', 'T', 'L', 'M', 'S', 'S', 'P', 0}
)

const invalid_MESSAGE_SIGNATURE = "invalid message signature"

const defaultFlags = nTLM_NEGOITATE_TARGET_INFO | nTLM_NEGOTIATE_56 | nTLM_NEGOTIATE_128 | nTLM_NEGOTIATE_UNICODE | nTLM_NEGOTIATE_EXTENDED_SESSION_SECURITY

type variableField struct {
	Length uint16
	MaxLen uint16
	Offset uint32
}

const data_EXTENDS_BEYOND_BOUNDARY = "data extends beyond boundary"

type version struct {
	Major    uint8
	Minor    uint8
	Build    uint16
	_        [3]byte // Reserved
	Revision uint8
}

type negotiate_Message struct {
	Signature   [8]byte
	MessageType uint32
	Flags       uint32
	Domain      variableField
	Workstation variableField
	Version     version
	// Payload is not in little endian, so we add it after the struct is converted to bytes.
}

const expectedNegotiateMessageSize = 40
const invalid_NEGOTIATE_MESSAGE_SIZE = "invalid negotiate message size"

type challenge_Header struct {
	Signature       [8]byte
	MessageType     uint32
	TargetName      variableField
	Flags           uint32
	ServerChallenge [8]byte
	Reserved        [8]byte // Reserved
	TargetInfo      variableField
}

type challenge_Message struct {
	Header        challenge_Header
	TargetName    string
	TargetInfo    map[uint16][]byte
	TargetInfoRaw []byte
}

const invalid_CHALLENGE_MESSAGE_TYPE = "invalid challenge message type"
const invalid_TARGET_INFO_LENGTH = "invalid target info length"

type authenticate_Message struct {
	Signature           [8]byte
	MessageType         uint32
	LmChallengeResponse variableField
	NtChallengeResponse variableField
	TargetName          variableField
	UserName            variableField
	Workstation         variableField
	_                   [8]byte // Reserved
	Flags               uint32
}

type authenticate_Crafter struct {
	LmChallengeResponse       []byte
	NtChallengeResponse       []byte
	TargetName                string
	UserName                  string
	EncryptedRandomSessionKey []byte
	Flags                     uint32
	MIC                       []byte
}

const no_ANONYMOUS_AUTH = "anonymous authentication is not supported"
const no_NTLMv1 = "NTLMv1 authentication is not supported"
const no_KEY_EXCH = "key exchange is not supported"
const must_UNICODE = "only unicode is supported"
