package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Key builds a stable cache key from method + model + request payload.
func Key(method, model string, payload any) string {
	b, err := json.Marshal(payload)
	if err != nil {
		b = []byte(fmt.Sprintf("%v", payload))
	}
	sum := sha256.Sum256(append([]byte(method+"|"+model+"|"), b...))
	return hex.EncodeToString(sum[:])
}
