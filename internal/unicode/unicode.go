package unicode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unicode/utf16"
)

const invalid_UNICODE_LENGTH = "invalid unicode length"

func FromUnicode(bin []byte) (string, error) {
	if len(bin)%2 != 0 {
		return "", errors.New(invalid_UNICODE_LENGTH)
	}
	uni := make([]uint16, len(bin)/2)
	if err := binary.Read(bytes.NewReader(bin), binary.LittleEndian, &uni); err != nil {
		return "", err
	}

	return string(utf16.Decode(uni)), nil
}

func ToUnicode(str string) []byte {
	uni := utf16.Encode([]rune(str))
	bin := bytes.Buffer{}
	binary.Write(&bin, binary.LittleEndian, &uni)
	return bin.Bytes()
}
