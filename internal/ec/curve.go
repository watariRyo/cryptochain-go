package ec

import (
	"crypto/elliptic"
	"math/big"
)

type CurveParams struct {
	P       *big.Int
	N       *big.Int
	B       *big.Int
	Gx, Gy  *big.Int
	BitSize int
	Name    string
}

// secp256k1の定義（https://www.secg.org/sec2-v2.pdf）
var (
	secp256k1 *CurveParams
	p, _      = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	n, _      = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	b, _      = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	gx, _     = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	gy, _     = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	bitSize   = 256
	name      = "secp256k1"
)

// 曲線のパラメータを返す
func (curve *CurveParams) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		P:       p,
		N:       n,
		B:       b,
		Gx:      gx,
		Gy:      gy,
		BitSize: bitSize,
		Name:    name,
	}
}

// 与えられた(x,y)が曲線上にあるかどうかを判定する
func (curve *CurveParams) IsOnCurve(x, y *big.Int) bool {
	y2 := new(big.Int).Exp(y, big.NewInt(2), nil)
	x3 := new(big.Int).Exp(x, big.NewInt(3), nil)
	ans := new(big.Int).Mod(y2.Sub(y2, x3.Add(x3, curve.B)), curve.P)

	return ans.Cmp(big.NewInt(0)) == 0
}

// (x1,y1) と (x2,y2) の和を返す
func (curve *CurveParams) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	z1 := zForAffine(x1, y1)
	z2 := zForAffine(x2, y2)
	return curve.affineFromJacobian(curve.addJacobian(x1, y1, z1, x2, y2, z2))
}

// 2*(x,y)を返す。
func (curve *CurveParams) Double(x1, y1 *big.Int) (x, y *big.Int) {
	z1 := zForAffine(x1, y1)
	return curve.affineFromJacobian(curve.doubleJacobian(x1, y1, z1))
}

// k*(Bx,By) を返す。k はビッグエンディアン形式の数値。
func (curve *CurveParams) ScalarMult(Bx, By *big.Int, k []byte) (x, y *big.Int) {
	Bz := new(big.Int).SetInt64(1)
	x, y, z := new(big.Int), new(big.Int), new(big.Int)

	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			x, y, z = curve.doubleJacobian(x, y, z) // 倍算
			if byte&0x80 == 0x80 {                  // ビットが1だったら加算
				x, y, z = curve.addJacobian(Bx, By, Bz, x, y, z)
			}
			byte <<= 1
		}
	}

	return curve.affineFromJacobian(x, y, z)
}

// k*G を返す
// kはビッグエンディアン形式の整数
func (curve *CurveParams) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return curve.ScalarMult(curve.Gx, curve.Gy, k)
}
