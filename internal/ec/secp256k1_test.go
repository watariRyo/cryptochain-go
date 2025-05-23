package ec

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
)

// ロジックではなくキー生成と曲線上の検証をする
func TestGenerateKey(t *testing.T) {
	// secp256k1曲線を取得
	curve := Secp256k1()

	// 秘密鍵を生成
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Errorf("failed to generate key: %v", err)
	}

	// 公開鍵を取得
	pubKey := privateKey.PublicKey
	t.Logf("Private Key: %x\n", privateKey.D.Bytes())
	t.Logf("Public Key (X): %x\n", pubKey.X.Bytes())
	t.Logf("Public Key (Y): %x\n", pubKey.Y.Bytes())

	// 曲線上の検証
	if !curve.IsOnCurve(pubKey.X, pubKey.Y) {
		t.Errorf("public key is not on curve")
	}

	// 16進数変換
	hex := PublicKeyToHexCompressed(pubKey.X, pubKey.Y)
	t.Logf("hex string: %s\n", hex)
}
