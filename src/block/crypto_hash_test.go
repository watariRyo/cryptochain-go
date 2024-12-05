package block

import (
	"testing"
)

func TestCryptoHashSHA256(t *testing.T) {
	expected := "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae"
	sum := cryptoHash("foo")

	if sum != expected {
		t.Errorf("cryptoHash failed to make hash. Expected %v, got %v", expected, sum)
	}

	expected = cryptoHash("one", "two", "three")
	sum = cryptoHash("three", "one", "two")

	if sum != expected {
		t.Errorf("cryptoHash mismatch. expected %v, got %v", expected, sum)
	}
}
