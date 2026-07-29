package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func newID(seed string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	s := hex.EncodeToString(b)
	seed = strings.TrimSpace(seed)
	if seed == "" {
		return s
	}
	seed = strings.ToLower(seed)
	seed = strings.ReplaceAll(seed, " ", "-")
	if len(seed) > 20 {
		seed = seed[:20]
	}
	return seed + "-" + s
}
