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

// 16進文字を2進数に変換する
func charToBinary(char rune) int {
	if char >= '0' && char <= '9' {
		return int(char - '0')
	} else if char >= 'a' && char <= 'f' {
		return int(char-'a') + 10
	} else if char >= 'A' && char <= 'F' {
		return int(char-'A') + 10
	}
	return 0 // 無効な文字の場合
}
