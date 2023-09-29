package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(payload []byte) string {
	sum := sha256.Sum256(payload)

	return hex.EncodeToString(sum[:])
}
