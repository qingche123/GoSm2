package sm2

import (
	"math/big"
)

type ECCurveParams struct {
	BitSize int
	P       *big.Int
	A       *big.Int
	B       *big.Int
	N       *big.Int
	Gx      *big.Int
	Gy      *big.Int
}

var G *ECPoint
var Infinity *ECPoint
var Ecurve *ECCurveParams

func InitSecp256k1() {

	Ecurve := &ECCurveParams{256, nil, nil, nil, nil, nil, nil}
	Ecurve.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	Ecurve.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	Ecurve.A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	Ecurve.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)

	Ecurve.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	Ecurve.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	Ecurve.BitSize = 256

	/*	X := &ECFieldElement{big.NewInt(0), Ecurve}
		Y := &ECFieldElement{big.NewInt(0), Ecurve}
		Infinity = &ECPoint{X, Y, Ecurve}*/
	Infinity = &ECPoint{nil, nil, Ecurve}

	GX := &ECFieldElement{Ecurve.Gx, Ecurve}
	GY := &ECFieldElement{Ecurve.Gy, Ecurve}
	G = &ECPoint{GX, GY, Ecurve}

	return
}

func InitSecp256r1() {

	Ecurve := &ECCurveParams{256, nil, nil, nil, nil, nil, nil}
	Ecurve.P, _ = new(big.Int).SetString("00FFFFFFFF00000001000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFF", 16)
	Ecurve.N, _ = new(big.Int).SetString("00FFFFFFFF00000001000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFC", 16)
	Ecurve.A, _ = new(big.Int).SetString("005AC635D8AA3A93E7B3EBBD55769886BC651D06B0CC53B0F63BCE3C3E27D2604B", 16)
	Ecurve.B, _ = new(big.Int).SetString("00FFFFFFFF00000000FFFFFFFFFFFFFFFFBCE6FAADA7179E84F3B9CAC2FC632551", 16)

	Ecurve.Gx, _ = new(big.Int).SetString("6B17D1F2E12C4247F8BCE6E563A440F277037D812DEB33A0F4A13945D898C296", 16)
	Ecurve.Gy, _ = new(big.Int).SetString("4FE342E2FE1A7F9B8EE7EB4A7C0F9E162BCE33576B315ECECBB6406837BF51F5", 16)
	Ecurve.BitSize = 256

	/*	X := &ECFieldElement{big.NewInt(0), Ecurve}
		Y := &ECFieldElement{big.NewInt(0), Ecurve}
		Infinity = &ECPoint{X, Y, Ecurve}*/
	Infinity = &ECPoint{nil, nil, Ecurve}

	GX := &ECFieldElement{Ecurve.Gx, Ecurve}
	GY := &ECFieldElement{Ecurve.Gy, Ecurve}
	G = &ECPoint{GX, GY, Ecurve}

	return
}
