package crypto

import (
	"encoding/json"
	"testing"
)

func TestCryptoHashSHA256(t *testing.T) {
	expected := "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"
	sum := CryptoHash("foo")

	if sum != expected {
		t.Errorf("cryptoHash failed to make hash. Expected %v, got %v", expected, sum)
	}

	expected = CryptoHash("one", "two", "three")
	sum = CryptoHash("three", "one", "two")

	if sum != expected {
		t.Errorf("cryptoHash mismatch. expected %v, got %v", expected, sum)
	}
}
func TestProducesUniqueHash(t *testing.T) {
	foo := make(map[string]int)
	bytes, _ := json.Marshal(foo)
	originalHash := CryptoHashByte(string(bytes))
	foo["a"] = 1
	newBytes, _ := json.Marshal(foo)
	nextHash := CryptoHashByte(string(newBytes))
	if string(originalHash) == string(nextHash) {
		t.Errorf("could not produces unique hash.")
	}
}

func TestCharToBinary(t *testing.T) {
	tests := []struct {
		input    rune
		expected int
	}{
		{'0', 0},
		{'1', 1},
		{'9', 9},
		{'a', 10},
		{'f', 15},
		{'A', 10},
		{'F', 15},
		{'z', 0}, // 無効な文字のケース
	}

	for _, tt := range tests {
		result := CharToBinary(tt.input)
		if result != tt.expected {
			t.Errorf("CharToBinary('%c') = %d; want %d", tt.input, result, tt.expected)
		}
	}
}
