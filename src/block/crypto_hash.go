package block

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

func cryptoHash(inputs ...string) string {
	sort.SliceStable(inputs, func(i, j int) bool {
		return inputs[i] < inputs[j]
	})
	result := strings.Join(inputs, "")
	hash := sha256.Sum256([]byte(result))
	return hex.EncodeToString(hash[:])
}
