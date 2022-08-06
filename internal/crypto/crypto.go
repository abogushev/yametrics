package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

func GetSign(key, msg string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))

	return fmt.Sprintf("%x", h.Sum(nil))
}
