package ec

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
)

var once sync.Once

func Secp256k1() elliptic.Curve {
	once.Do(func() {
		secp256k1 = &CurveParams{
			P:       p,
			N:       n,
			B:       b,
			Gx:      gx,
			Gy:      gy,
			BitSize: bitSize,
			Name:    name,
		}
	})
	return secp256k1
}

// 公開鍵の圧縮
func compressPublicKey(x, y *big.Int) []byte {
	// X座標をバイトスライスに変換
	xBytes := x.Bytes()

	// プレフィックスの設定
	var prefix byte
	if y.Bit(0) == 0 { // Y座標が偶数
		prefix = 0x02
	} else { // Y座標が奇数
		prefix = 0x03
	}

	// 圧縮形式を返す
	return append([]byte{prefix}, xBytes...)
}

// 圧縮形式で公開鍵を16進数に変換する
func PublicKeyToHexCompressed(x, y *big.Int) string {
	compressed := compressPublicKey(x, y)
	return hex.EncodeToString(compressed)
}

// 復元
func DecompressHexPublicKey(curve elliptic.Curve, compressedHex string) (*big.Int, *big.Int, error) {
	compressed, err := hex.DecodeString(compressedHex)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid hex public key ")
	}

	if len(compressed) != 33 {
		return nil, nil, fmt.Errorf("invalid compressed public key length")
	}

	// プレフィックスを取得
	prefix := compressed[0]
	if prefix != 0x02 && prefix != 0x03 {
		return nil, nil, fmt.Errorf("invalid prefix for compressed public key")
	}

	// X座標を取得
	x := new(big.Int).SetBytes(compressed[1:])

	// 曲線のパラメータを取得
	p := curve.Params().P

	// Y^2 = X^3 + B (mod P) を計算
	y2 := new(big.Int).Mul(x, x)
	y2.Mul(y2, x)
	y2.Add(y2, curve.Params().B)
	y2.Mod(y2, p)

	// Y座標を復元
	y := new(big.Int).ModSqrt(y2, p)
	if y == nil {
		return nil, nil, fmt.Errorf("failed to compute square root of Y^2")
	}

	// Y座標の偶奇を確認し、必要に応じて反転
	if (y.Bit(0) == 1 && prefix == 0x02) || (y.Bit(0) == 0 && prefix == 0x03) {
		y.Sub(p, y)
	}

	return x, y, nil
}
