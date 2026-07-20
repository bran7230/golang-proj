package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// verify is used in the decode method to check the HMAC signature in the request
// header
func verify(msg, key []byte, hash string) (bool, error) {
	sig, err := hex.DecodeString(hash)

	if err != nil {
		return false, err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(msg)

	return hmac.Equal(sig, mac.Sum(nil)), nil
}
