package ntlm

import (
	"crypto/hmac"
	"crypto/md5"
	"strings"

	"github.com/huner2/go-wsus/internal/unicode"
	//lint:ignore SA1019 Need to use md4 for NTLMv2
	"golang.org/x/crypto/md4"
)

func generateHash(user, pass, target string) []byte {
	return hmacMd5(generateNTLM(pass), unicode.ToUnicode(strings.ToUpper(user)+target))
}

func generateNTLM(pass string) []byte {
	hash := md4.New()
	hash.Write(unicode.ToUnicode(pass))
	return hash.Sum(nil)
}

func hmacMd5(key []byte, data ...[]byte) []byte {
	mac := hmac.New(md5.New, key)
	for _, d := range data {
		mac.Write(d)
	}
	return mac.Sum(nil)
}

func computeNTLMv2(hash, serverChall, clientChall, timestamp, targetInfo []byte) []byte {
	temp := []byte{1, 1, 0, 0, 0, 0, 0, 0}
	temp = append(temp, timestamp...)
	temp = append(temp, clientChall...)
	temp = append(temp, 0, 0, 0, 0)
	temp = append(temp, targetInfo...)
	temp = append(temp, 0, 0, 0, 0)

	proofStr := hmacMd5(hash, serverChall, temp)
	return append(proofStr, temp...)
}

func computeLMv2(hash, serverChall, clientChall []byte) []byte {
	return append(hmacMd5(hash, serverChall, clientChall), clientChall...)
}
