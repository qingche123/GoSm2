package sm2

import (
	"fmt"
	"math/big"
)

//----------------------------------------------------------------------

type ECPoint struct {
	X, Y  *ECFieldElement
	curve *ECCurveParams
}

func NewECPoint() *ECPoint {
	/*
		bnx := big.NewInt(0)
		bny := big.NewInt(0)

		x := &ECFieldElement{bnx, Ecurve}
		y := &ECFieldElement{bny, Ecurve}
	*/

	x := NewECFieldElement()
	y := NewECFieldElement()

	return &ECPoint{x, y, Ecurve}
}

func DumpECPoint(dst *ECPoint, src *ECPoint) {

	DumpECFieldElement(dst.X, src.X)
	DumpECFieldElement(dst.Y, src.Y)

	dst.curve = src.curve
}

func (e *ECPoint) IsInfinity() bool {
	if e.X == nil && e.Y == nil {
		return false
	} else {
		return true
	}
}

func (e *ECPoint) GetSize() int {
	if e.IsInfinity() {
		return 1
	} else {
		return 33
	}
}

func (e *ECPoint) CompareTo(other *ECPoint) int {
	if e.X.value.Cmp(other.X.value) == 0 && e.Y.value.Cmp(other.Y.value) == 0 {
		return 0
	} else {
		return 1
	}
}

func DecompressPoint(yTilde int64, X1 *big.Int, curve *ECCurveParams) *ECPoint {
	x := &ECFieldElement{X1, curve}
	y := x.Square()
	y.AddBig(y, curve.A)
	x.Mul(y, x)
	x.AddBig(x, curve.B)
	beta := x.Sqrt()

	if beta == nil {
		fmt.Println("Invalid point compression")
	}

	betaValue := new(big.Int).Set(beta.value)
	bit0 := new(big.Int).Mod(betaValue, big.NewInt(2))
	if bit0.Int64() != yTilde {
		beta = &ECFieldElement{big.NewInt(0).Sub(curve.P, betaValue), curve}
	}

	return &ECPoint{x, beta, curve}
}

func DecodePoint(encoded []byte, curve *ECCurveParams) *ECPoint {
	var p *ECPoint = nil

	expectedLength := (curve.P.BitLen() + 7) / 8

	switch encoded[0] {
	case 0x00:
		if len(encoded) != 1 {
			fmt.Println("Incorrect length for infinity encoding", "encoded")
		}
		p = Infinity
	case 0x02, 0x03:
		{
			if len(encoded) != expectedLength+1 {
				fmt.Println("Incorrect length for infinity encoding", "encoded")
			}
			yTilde := encoded[0] & 1
			Tmp := encoded[1:]
			Reverse(Tmp)
			Tmp1 := make([]byte, len(Tmp)+1)
			copy(Tmp1, Tmp)

			X1 := new(big.Int).SetBytes(Tmp1)
			p = DecompressPoint(int64(yTilde), X1, curve)
		}
	case 0x04, 0x06, 0x07:
		{
			if len(encoded) != (2*expectedLength + 1) {
				fmt.Println("Incorrect length for uncompressed/hybrid encoding", "encoded")
			}

			Tmp1 := encoded[1 : 1+expectedLength]
			Reverse(Tmp1)
			Tmp2 := make([]byte, len(Tmp1)+1)
			copy(Tmp2, Tmp1)
			X1 := new(big.Int).SetBytes(Tmp2)

			Tmp3 := encoded[1:]
			Reverse(Tmp3)
			Tmp4 := make([]byte, len(Tmp3)+1)
			copy(Tmp4, Tmp3)
			Y1 := new(big.Int).SetBytes(Tmp4)

			Ex := &ECFieldElement{X1, curve}
			Ey := &ECFieldElement{Y1, curve}
			p = &ECPoint{Ex, Ey, curve}
		}
	default:
		fmt.Println("Invalid point encoding ", encoded[0])
	}
	return p
}

func (e *ECPoint) EncodePoint(compressed bool) []byte {
	if e.IsInfinity() {
		Tmp := make([]byte, 1)
		return Tmp
	}
	var data []byte
	if compressed {
		data = make([]byte, 33)
	} else {
		data = make([]byte, 65)
		yBytes := e.Y.value.Bytes()
		Reverse(yBytes)
		copy(data[65-len(data):], yBytes)
	}

	xBytes := e.X.value.Bytes()
	Reverse(xBytes)
	copy(data[33-len(data):], xBytes)

	if !compressed {
		data[0] = 0x04
	} else {
		Tmp := big.NewInt(0)
		Tmp.Mod(e.Y.value, big.NewInt(2))
		if Tmp.Int64() == 0 {
			data[0] = 0x02
		} else {
			data[0] = 0x03
		}
	}
	return data
}

func (e *ECPoint) Equals(other *ECPoint) bool {
	if e == other {
		return true
	}
	if e.IsInfinity() && other.IsInfinity() {
		return true
	}
	if e.IsInfinity() || other.IsInfinity() {
		return false
	}

	if e.CompareTo(other) == 0 {
		return true
	} else {
		return false
	}
}

func (e *ECPoint) Twice() *ECPoint {
	if e.IsInfinity() {
		return e
	}
	if e.Y.value.Sign() == 0 {
		return Infinity
	}

	TWO := &ECFieldElement{big.NewInt(2), e.curve}
	THREE := &ECFieldElement{big.NewInt(3), e.curve}

	Tmp1 := big.NewInt(0)
	Tmp1.Exp(e.X.value, big.NewInt(2), big.NewInt(0))

	Tmp2 := &ECFieldElement{big.NewInt(0), e.curve}
	Tmp2.MulBig(THREE, Tmp1)
	Tmp2.AddBig(Tmp2, e.curve.A)

	Tmp3 := TWO.MulBig(TWO, e.Y.value)

	gamma := Tmp2.Div(Tmp2, Tmp3)

	Tmp4 := gamma.Square()
	Tmp5 := TWO.MulBig(TWO, e.X.value)
	x3 := Tmp4.Sub(Tmp4, Tmp5)

	y3 := &ECFieldElement{big.NewInt(0), e.curve}
	y3.Sub(e.X, x3)
	y3.Mul(y3, gamma)
	y3.Sub(y3, e.Y)

	return &ECPoint{x3, y3, e.curve}
}

func Multiply(p *ECPoint, k *big.Int) *ECPoint {
	m := k.BitLen()

	var width byte
	var reqPreCompLen int

	if m < 13 {
		width = 2
		reqPreCompLen = 1
	} else if m < 41 {
		width = 3
		reqPreCompLen = 2
	} else if m < 121 {
		width = 4
		reqPreCompLen = 4
	} else if m < 337 {
		width = 5
		reqPreCompLen = 8
	} else if m < 897 {
		width = 6
		reqPreCompLen = 16
	} else if m < 2305 {
		width = 7
		reqPreCompLen = 32
	} else {
		width = 8
		reqPreCompLen = 127
	}
	preCompLen := 1

	preComp := make([]ECPoint, reqPreCompLen)
	twiceP := p.Twice()

	if preCompLen < reqPreCompLen {
		oldPreComp := preComp
		preComp := make([]ECPoint, reqPreCompLen)
		copy(preComp, oldPreComp)

		i := preCompLen

		for i < reqPreCompLen {
			preComp[i].Add(twiceP, &preComp[i-1])
			i++
		}
	}

	wnaf := WindowNaf(width, k)
	l := len(wnaf)

	q := Infinity
	i := l - 1
	for i >= 0 {
		q = q.Twice()
		if wnaf[i] != 0 {
			if wnaf[i] > 0 {

				q.Add(q, &preComp[(wnaf[i]-1)/2])
			}
		}
		i--
	}

	return q

}

func WindowNaf(width byte, k *big.Int) []byte {
	wnaf := make([]byte, k.BitLen()+1)
	var pow2wB uint16

	pow2wB = 1 << width
	i := 0
	length := 0
	bigp2wB := big.NewInt(int64(pow2wB))
	Tmp := big.NewInt(0)

	for k.Sign() > 0 {
		if !IsEven(k) {
			remainder := big.NewInt(0)
			remainder.Mod(k, bigp2wB)
			if remainder.Bit(int(width-1)) == 1 {
				Tmp.Sub(remainder, bigp2wB)
				wnaf[i] = byte(Tmp.Int64())
			} else {
				wnaf[i] = byte(remainder.Int64())
			}
			k.Sub(k, big.NewInt(int64(wnaf[i])))
			length = i
		} else {
			wnaf[i] = 0
		}
		k.Rsh(k, 1)
		i++
	}
	length++
	wnafShort := make([]byte, length)
	copy(wnafShort, wnaf)
	return wnafShort
}

func (e *ECPoint) Neg(x *ECPoint) *ECPoint {
	return &ECPoint{x.X, x.Y.Neg(x.Y), x.curve}
}

func (e *ECPoint) Mul(p *ECPoint, n []byte) *ECPoint {
	if nil == p || nil == n {
		fmt.Println("Argument is Null")
	}
	if len(n) != 32 {
		fmt.Println("Argument is not excepted")
	}
	if p.IsInfinity() {
		return p
	}
	k := big.NewInt(0)
	l := len(n)
	Tmp1 := make([]byte, l)

	copy(Tmp1, n)
	Reverse(Tmp1)
	Tmp2 := make([]byte, l+1)
	copy(Tmp2, n)

	k.SetBytes(Tmp2)
	if k.Sign() == 0 {
		return Infinity
	}
	return Multiply(p, k)
}

func (e *ECPoint) Add(x *ECPoint, y *ECPoint) *ECPoint {
	if x.IsInfinity() {
		return y
	}
	if y.IsInfinity() {
		return x
	}
	if x.X.Equals(y.X) {
		if x.Y.Equals(y.Y) {
			return x.Twice()
		}
		return Infinity
	}
	Tmp1 := new(ECFieldElement)
	Tmp1.Sub(y.Y, x.Y)
	Tmp2 := new(ECFieldElement)
	Tmp2.Sub(y.X, x.X)
	gama := Tmp1.Div(Tmp1, Tmp2)

	x3 := gama.Square()
	x3.Sub(x3, x.X)
	x3.Sub(x3, y.X)

	y3 := x.X.Sub(x.X, x3)
	y3.Mul(y3, gama)
	y3.Sub(y3, x.Y)

	return &ECPoint{x3, y3, x.curve}

}

func (e *ECPoint) Sub(x *ECPoint, y *ECPoint) *ECPoint {
	if y.IsInfinity() {
		return x
	}
	return x.Add(x, y.Neg(y))
}
