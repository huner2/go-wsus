// Much of the structure of the code in this module is based on the Azure ntlmssp repository: https://github.com/Azure/go-ntlmssp
// Azure ntlmssp is licensed under the MIT license.

package ntlm

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/huner2/go-wsus/internal/unicode"
)

// Create fields needed for length, as well as determine offset by given pointer
func newVariableField(offset *int, size int) variableField {
	defer func() { *offset += size }()
	return variableField{
		Length: uint16(size),
		MaxLen: uint16(size),
		Offset: uint32(*offset),
	}
}

func (v *variableField) unmarshal(data []byte) ([]byte, error) {
	if len(data) < int(v.Offset+uint32(v.Length)) {
		return nil, errors.New(data_EXTENDS_BEYOND_BOUNDARY)
	}
	return data[v.Offset : v.Offset+uint32(v.Length)], nil
}

func (v *variableField) unmarshalString(data []byte, uni bool) (string, error) {
	read, err := v.unmarshal(data)
	if err != nil {
		return "", err
	}
	if uni {
		return unicode.FromUnicode(read)
	}
	// OEM encoding
	return string(read), nil
}

func validateSignature(signature [8]byte) bool {
	return bytes.Equal(signature[:], nTLM_MESSAGE_SIGNATURE[:])
}

func (n *NTLMNegotiator) newNegotiateMessage(domain, workstation string) ([]byte, error) {
	offset := expectedNegotiateMessageSize
	flags := defaultFlags

	if domain != "" {
		flags |= nTLM_NEGOTIATE_DOMAIN_SUPPLIED
	}
	if workstation != "" {
		flags |= nTLM_NEGOTIATE_WORKSTATION_SUPPLIED
	}

	msg := negotiate_Message{
		Signature:   nTLM_MESSAGE_SIGNATURE,
		MessageType: 1,
		Flags:       flags,
		Domain:      newVariableField(&offset, len(domain)),
		Workstation: newVariableField(&offset, len(workstation)),
		Version:     version{Major: 6, Minor: 1, Build: 7601, Revision: 15},
	}

	buf := bytes.Buffer{}
	if err := binary.Write(&buf, binary.LittleEndian, &msg); err != nil {
		if n.Debug {
			log.Printf("[DEBUG] [Marshaling negotiate message] %s", err)
		}
		return nil, err
	}
	if buf.Len() != expectedNegotiateMessageSize {
		return nil, errors.New(invalid_NEGOTIATE_MESSAGE_SIZE)
	}

	// Add domain and workstation strings to the end of the buffer
	// Note: The strings are not in little endian, so we add them after the struct is converted to bytes.
	payload := strings.ToUpper(domain + workstation)
	if _, err := buf.WriteString(payload); err != nil {
		if n.Debug {
			log.Printf("[DEBUG] [Writing domain and workstation strings] %s", err)
		}
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *challenge_Message) unmarshal(data []byte, debug bool) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &c.Header.Signature); err != nil {
		return err
	}
	if !validateSignature(c.Header.Signature) {
		if debug {
			log.Printf("[DEBUG] Invalid message signature, got %s\n", c.Header.Signature)
		}
		return errors.New(invalid_MESSAGE_SIGNATURE)
	}
	if err := binary.Read(reader, binary.LittleEndian, &c.Header.MessageType); err != nil {
		return err
	}
	if c.Header.MessageType != 2 {
		if debug {
			log.Printf("[DEBUG] Invalid message type, got %d\n", c.Header.MessageType)
		}
		return errors.New(invalid_CHALLENGE_MESSAGE_TYPE)
	}
	if err := binary.Read(reader, binary.LittleEndian, &c.Header.TargetName); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &c.Header.Flags); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &c.Header.ServerChallenge); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &c.Header.Reserved); err != nil {
		return err
	}
	if err := binary.Read(reader, binary.LittleEndian, &c.Header.TargetInfo); err != nil {
		return err
	}

	if c.Header.TargetName.Length > 0 {
		var err error
		c.TargetName, err = c.Header.TargetName.unmarshalString(data, c.Header.Flags&nTLM_NEGOTIATE_UNICODE == nTLM_NEGOTIATE_UNICODE)
		if err != nil {
			return err
		}
	}

	if c.Header.TargetInfo.Length > 0 {
		bin, err := c.Header.TargetInfo.unmarshal(data)
		if err != nil {
			return err
		}
		c.TargetInfoRaw = bin
		c.TargetInfo = make(map[uint16][]byte)
		reader = bytes.NewReader(bin)
		for {
			var id uint16
			var len uint16
			if err := binary.Read(reader, binary.LittleEndian, &id); err != nil {
				return err
			}
			if id == avID_MsvAvEOL {
				break
			}

			if err := binary.Read(reader, binary.LittleEndian, &len); err != nil {
				return err
			}
			value := make([]byte, len)
			n, err := reader.Read(value)
			if err != nil {
				return err
			}
			if n != int(len) {
				return errors.New(invalid_TARGET_INFO_LENGTH)
			}
			c.TargetInfo[id] = value
		}
	}

	return nil
}

func (n *NTLMNegotiator) craftResponse(challenge []byte, user, password string, isHash bool) ([]byte, error) {
	if user == "" && password == "" {
		return nil, errors.New(no_ANONYMOUS_AUTH)
	}

	var challengeMessage challenge_Message
	if err := challengeMessage.unmarshal(challenge, n.Debug); err != nil {
		if n.Debug {
			log.Printf("[DEBUG] [Unmarshaling challenge message] %s", err)
		}
		return nil, err
	}

	if challengeMessage.Header.Flags&nTLM_NEGOTIATE_LM_KEY == nTLM_NEGOTIATE_LM_KEY {
		return nil, errors.New(no_NTLMv1)
	}
	if challengeMessage.Header.Flags&nTLM_NEGOTIATE_KEY_EXCH == nTLM_NEGOTIATE_KEY_EXCH {
		return nil, errors.New(no_KEY_EXCH)
	}

	authMessage := authenticate_Crafter{
		UserName:   user,
		TargetName: challengeMessage.TargetName,
		Flags:      challengeMessage.Header.Flags,
	}

	timestamp := challengeMessage.TargetInfo[avID_MsvAvTimestamp]
	if timestamp == nil {
		ft := uint64(time.Now().UnixNano()) / 100
		ft += 116444736000000000 // add time between unix & windows offset
		timestamp = make([]byte, 8)
		binary.LittleEndian.PutUint64(timestamp, ft)
	}

	clientChallenge := make([]byte, 8)
	rand.Reader.Read(clientChallenge)

	var hash []byte
	if isHash {
		var temp string
		var err error
		hashParts := strings.Split(password, ":")
		if len(hashParts) > 1 {
			temp = hashParts[1]
		}
		hash, err = hex.DecodeString(temp)
		if err != nil {
			if n.Debug {
				log.Printf("[DEBUG] [Decoding hash] %s", err)
			}
			return nil, err
		}
		hash = hmacMd5(hash, unicode.ToUnicode(strings.ToUpper(user)+challengeMessage.TargetName))
	} else {
		hash = generateHash(user, password, challengeMessage.TargetName)
	}

	authMessage.NtChallengeResponse = computeNTLMv2(hash, challengeMessage.Header.ServerChallenge[:], clientChallenge, timestamp, challengeMessage.TargetInfoRaw)

	if challengeMessage.TargetInfoRaw == nil {
		authMessage.LmChallengeResponse = computeLMv2(hash, challengeMessage.Header.ServerChallenge[:], clientChallenge)
	}

	data, err := authMessage.marshal()
	if err != nil {
		if n.Debug {
			log.Printf("[DEBUG] [Marshaling authenticate message] %s", err)
		}
		return nil, err
	}
	return data, nil
}

func (a *authenticate_Crafter) marshal() ([]byte, error) {
	if a.Flags&nTLM_NEGOTIATE_UNICODE != nTLM_NEGOTIATE_UNICODE {
		return nil, errors.New(must_UNICODE)
	}

	target, user := unicode.ToUnicode(a.TargetName), unicode.ToUnicode(a.UserName)
	workstation := unicode.ToUnicode("")

	ptr := binary.Size(&authenticate_Message{})
	message := authenticate_Message{
		Signature:           nTLM_MESSAGE_SIGNATURE,
		MessageType:         3,
		LmChallengeResponse: newVariableField(&ptr, len(a.LmChallengeResponse)),
		NtChallengeResponse: newVariableField(&ptr, len(a.NtChallengeResponse)),
		TargetName:          newVariableField(&ptr, len(target)),
		UserName:            newVariableField(&ptr, len(user)),
		Workstation:         newVariableField(&ptr, len(workstation)),
	}

	a.Flags &= ^nTLM_NEGOTIATE_VERSION

	bin := bytes.Buffer{}
	if err := binary.Write(&bin, binary.LittleEndian, &message); err != nil {
		return nil, err
	}
	if err := binary.Write(&bin, binary.LittleEndian, a.LmChallengeResponse); err != nil {
		return nil, err
	}
	if err := binary.Write(&bin, binary.LittleEndian, a.NtChallengeResponse); err != nil {
		return nil, err
	}
	if err := binary.Write(&bin, binary.LittleEndian, target); err != nil {
		return nil, err
	}
	if err := binary.Write(&bin, binary.LittleEndian, user); err != nil {
		return nil, err
	}
	if err := binary.Write(&bin, binary.LittleEndian, workstation); err != nil {
		return nil, err
	}

	return bin.Bytes(), nil
}
