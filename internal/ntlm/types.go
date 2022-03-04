package ntlm

const (
	NTLM_NEGOTIATE_UNICODE = 1 << iota
	NTLM_NEGOTIATE_OEM     = 1 << iota
	NTLM_REQUEST_TARGET    = 1 << iota
	_
	NTLM_NEGOTIATE_SIGN     = 1 << iota
	NTLM_NEGOTIATE_SEAL     = 1 << iota
	NTLM_NEGOTIATE_DATAGRAM = 1 << iota
	NTLM_NEGOTIATE_LM_KEY   = 1 << iota
	_
	NTLM_NEGOTIATE_NTLM = 1 << iota
	_
	NTLM_NEGOTIATE_ANONYMOUS            = 1 << iota
	NTLM_NEGOTIATE_DOMAIN_SUPPLIED      = 1 << iota
	NTLM_NEGOTIATE_WORKSTATION_SUPPLIED = 1 << iota
	_
	NTLM_NEGOTIATE_ALWAYS_SIGN = 1 << iota
	NTLM_TARGET_TYPE_DOMAIN    = 1 << iota
	NTLM_TARGET_TYPE_SERVER    = 1 << iota
	_
	NTLM_NEGOTIATE_EXTENDED_SESSION_SECURITY = 1 << iota
	NTLM_NEGOTIATE_IDENTIFY                  = 1 << iota
	_
	NTLM_NEGOTIATE_REQUEST_NON_NT_KEY = 1 << iota
	NTLM_NEGOITATE_TARGET_INFO        = 1 << iota
	_
	NTLM_NEGOTIATE_VERSION = 1 << iota
	_
	_
	_
	NTLM_NEGOTIATE_128      = 1 << iota
	NTLM_NEGOTIATE_KEY_EXCH = 1 << iota
	NTLM_NEGOTIATE_56       = 1 << iota
)

type Negotiate_Message struct {
	Signature   [8]byte
	MessageType [4]byte
	Flags       [4]byte
}

type Authenticate_Message struct {
	Signature         [8]byte
	Type              [4]byte
	LmResponseLen     [2]byte
	LmResponseMaxLen  [2]byte
	LmResponseOffset  [4]byte
	NtResponseLen     [2]byte
	NtResponseMaxLen  [2]byte
	NtResponseOffset  [4]byte
	DomainLen         [2]byte
	DomainMaxLen      [2]byte
	DomainOffset      [4]byte
	UserLen           [2]byte
	UserMaxLen        [2]byte
	UserOffset        [4]byte
	WorkstationLen    [2]byte
	WorkstationMaxLen [2]byte
	WorkstationOffset [4]byte
	SessionKeyLen     [2]byte
	SessionKeyMaxLen  [2]byte
	SessionKeyOffset  [4]byte
	Flags             [4]byte
	Version           [8]byte
}
