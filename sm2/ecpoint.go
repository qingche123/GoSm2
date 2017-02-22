package sm2

import (
	"fmt"
	"math/big"
)

// ECPoint
type ECPoint struct {
	X, Y  *ECFieldElement
	curve *ECCurveParams
}

func PrintHex(str string, bt []byte, length int) {
	fmt.Println(str, "Length = ", length)
	for i := 0; i < length; i++ {
		if i%16 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("0x%02x ", bt[i])
	}
	fmt.Println(" ")
	fmt.Println(" ")
}

func PrintHexEx(str string, bt []byte, length int) {
	fmt.Println(str, "Length = ", length)
	for i := 0; i < length; i++ {
		if i%16 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("%02x", bt[i])
	}
	fmt.Println(" ")
	fmt.Println(" ")
}

func NewECPoint() *ECPoint {

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
		return true
	}
	return false
}

func (e *ECPoint) GetSize() int {
	if e.IsInfinity() {
		return 1
	}
	return 33

}

func (e *ECPoint) CompareTo(other *ECPoint) int {
	if e.X.value.Cmp(other.X.value) == 0 && e.Y.value.Cmp(other.Y.value) == 0 {
		return 0
	}
	return 1
}

func DecompressPoint(yTilde int64, X1 *big.Int, curve *ECCurveParams) *ECPoint {
	x := &ECFieldElement{X1, curve}
	y := x.Square()
	y.AddBig(y, curve.A)
	x.Mul(y, x)
	x.AddBig(x, curve.B)

	beta := x.Sqrt()

	//PrintHex("beta====", beta.value.Bytes(), len(beta.value.Bytes()))
	if beta == nil {
		fmt.Println("Invalid point compression")
		//os.Exit(0)
	}

	betaValue := new(big.Int).Set(beta.value)

	bit0 := big.NewInt(0)
	bit0.Mod(betaValue, big.NewInt(2))

	if bit0.Int64() != yTilde {
		v := big.NewInt(0)
		v.Sub(curve.P, betaValue)
		beta = &ECFieldElement{v, curve}
	}

	return &ECPoint{x, beta, curve}
}

func DecodePoint(encoded []byte, curve *ECCurveParams) *ECPoint {

	p := NewECPoint()

	expectedLength := (curve.P.BitLen() + 7) / 8

	switch encoded[0] {
	case 0x00: // infinity
		if len(encoded) != 1 {
			fmt.Println("Incorrect length for infinity encoding", "encoded")
		}
		DumpECPoint(p, Infinity)
	case 0x02, 0x03: // compressed
		{
			if len(encoded) != expectedLength+1 {
				fmt.Println("Incorrect length for infinity encoding", "encoded")
			}
			yTilde := encoded[0] & 1

			tmp := make([]byte, len(encoded)-1)
			copy(tmp, encoded[1:])
			Reverse(tmp)
			PrintHex("xBytes", tmp, len(tmp))

			X1 := new(big.Int).SetBytes(tmp)

			p = DecompressPoint(int64(yTilde), X1, curve)
		}
	case 0x04, 0x06, 0x07: // uncompressed, hybrid, hybrid
		{
			if len(encoded) != (2*expectedLength + 1) {
				fmt.Println("Incorrect length for uncompressed/hybrid encoding", "encoded")
			}

			tmp1 := encoded[1 : 1+expectedLength]
			tmp2 := make([]byte, len(tmp1))
			copy(tmp2, tmp1)
			Reverse(tmp2)

			PrintHex("tmp2", tmp2, len(tmp2))
			X1 := new(big.Int).SetBytes(tmp2)

			tmp3 := encoded[1+expectedLength:]
			tmp4 := make([]byte, len(tmp3))
			copy(tmp4, tmp3)
			Reverse(tmp4)

			PrintHex("tmp4", tmp4, len(tmp4))
			Y1 := new(big.Int).SetBytes(tmp4)

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
		tmp := make([]byte, 1)
		fmt.Println("IsInfinity")
		return tmp
	}

	var data []byte

	if compressed {
		data = make([]byte, 33)
	} else {
		data = make([]byte, 65)

		yBytes := e.Y.value.Bytes()

		tmp := make([]byte, len(yBytes))
		copy(tmp, yBytes)
		Reverse(tmp)
		copy(data[65-len(yBytes):], tmp)
	}

	xBytes := e.X.value.Bytes()
	PrintHex("xBytes", xBytes, len(xBytes))

	tmp := make([]byte, len(xBytes))
	copy(tmp, xBytes)
	Reverse(tmp)

	copy(data[33-len(tmp):], tmp)

	if !compressed {
		data[0] = 0x04
	} else {
		if IsEven(e.Y.value) {
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
	}
	return false
}

// Twice ---
func (e *ECPoint) Twice() *ECPoint {
	if e.IsInfinity() {
		return e
	}
	if e.Y.value.Sign() == 0 {
		return e
	}

	TWO := &ECFieldElement{big.NewInt(2), e.curve}
	THREE := &ECFieldElement{big.NewInt(3), e.curve}

	Tmp1 := big.NewInt(0)
	Tmp1.Exp(e.X.value, big.NewInt(2), big.NewInt(0))

	Tmp2 := &ECFieldElement{big.NewInt(0), e.curve}
	Tmp2.MulBig(THREE, Tmp1)
	Tmp2.AddBig(Tmp2, e.curve.A)

	Tmp3 := NewECFieldElement()
	Tmp3.MulBig(TWO, e.Y.value)

	gamma := Tmp2.Div(Tmp2, Tmp3)

	Tmp4 := gamma.Square()
	Tmp5 := NewECFieldElement()
	Tmp5.MulBig(TWO, e.X.value)
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
	fmt.Println("reqPreCompLen = ", reqPreCompLen)

	preComp := make([]*ECPoint, reqPreCompLen)
	for i := 0; i < reqPreCompLen; i++ {
		preComp[i] = NewECPoint()
	}

	DumpECPoint(preComp[0], p)

	// vpX := p.X.value.Bytes()
	// vpY := p.Y.value.Bytes()
	// PrintHex("p point", vpX, len(vpX))
	// PrintHex("p point", vpY, len(vpY))

	twiceP := p.Twice()

	// vX := twiceP.X.value.Bytes()
	// vY := twiceP.Y.value.Bytes()
	// PrintHex("twiceP point", vX, len(vX))
	// PrintHex("twiceP point", vY, len(vY))

	if preCompLen < reqPreCompLen {
		oldPreComp := NewECPoint()
		DumpECPoint(oldPreComp, preComp[0])
		DumpECPoint(preComp[0], oldPreComp)

		for i := preCompLen; i < reqPreCompLen; i++ {
			preComp[i].Add(twiceP, preComp[i-1])
		}
	}
	// fmt.Println("width:=====", width)
	// PrintHex("k", k.Bytes(), len(k.Bytes()))

	wnaf := WindowNaf(width, k)
	l := len(wnaf)
	//PrintHex("wnaf", wnaf, l)
	//q := &ECPoint(nil, nil, Ecurve)

	//q := &ECPoint{nil, nil, Ecurve}
	//q := NewECPoint()
	//DumpECPoint(q, Infinity)
	q := NewECPoint()
	for i := l - 1; i >= 0; i-- {
		q = q.Twice()
		id := wnaf[i]
		if id != 0 {
			if id > 0 {
				index := (id - 1) / 2
				//fmt.Println("i：===", i, "index========", index)

				// fmt.Println("----------------------------------")
				// PrintHex("preComp[index]X", preComp[index].X.value.Bytes(), len(preComp[index].X.value.Bytes()))
				// fmt.Println("----------------------------------")
				// PrintHex("preComp[index]Y", preComp[index].Y.value.Bytes(), len(preComp[index].Y.value.Bytes()))

				q.Add(q, preComp[index])

				// fmt.Println("----------------------------------")
				// PrintHex("pbkeyX", q.X.value.Bytes(), len(q.X.value.Bytes()))
				// fmt.Println("----------------------------------")
				// PrintHex("pbkeyY", q.Y.value.Bytes(), len(q.Y.value.Bytes()))
				// os.Exit(0)
			} else {
				index := (-id - 1) / 2
				// fmt.Println("i：===", i, "index========", index)

				// fmt.Println("----------------------------------")
				// PrintHex("pbkeyX", q.X.value.Bytes(), len(q.X.value.Bytes()))
				// fmt.Println("----------------------------------")
				// PrintHex("pbkeyY", q.Y.value.Bytes(), len(q.Y.value.Bytes()))

				q.Sub(q, preComp[index])

				// fmt.Println("----------------------------------")
				// PrintHex("pbkeyX", q.X.value.Bytes(), len(q.X.value.Bytes()))
				// fmt.Println("----------------------------------")
				// PrintHex("pbkeyY", q.Y.value.Bytes(), len(q.Y.value.Bytes()))
				// os.Exit(0)
			}
		}
		// fmt.Println("----------------------------------")
		// PrintHex("pbkeyX", q.X.value.Bytes(), len(q.X.value.Bytes()))
		// fmt.Println("----------------------------------")
		// PrintHex("pbkeyY", q.Y.value.Bytes(), len(q.Y.value.Bytes()))
	}
	fmt.Println("----------------------------------")
	PrintHex("pbkeyX", q.X.value.Bytes(), len(q.X.value.Bytes()))
	fmt.Println("----------------------------------")
	PrintHex("pbkeyY", q.Y.value.Bytes(), len(q.Y.value.Bytes()))
	return q

}

// WindowNaf ---
func WindowNaf(width byte, k *big.Int) []int8 {
	wnaf := make([]int8, k.BitLen()+1)
	var pow2wB uint16

	pow2wB = 1 << width

	length := 0
	bigp2wB := big.NewInt(int64(pow2wB))
	Tmp := big.NewInt(0)

	for i := 0; k.Sign() > 0; i++ {
		if !IsEven(k) {
			remainder := big.NewInt(0)
			remainder.Mod(k, bigp2wB)
			if remainder.Bit(int(width-1)) == 1 {
				Tmp.Sub(remainder, bigp2wB)

				wnaf[i] = int8(Tmp.Int64())

			} else {
				wnaf[i] = int8(remainder.Int64())

			}
			k.Sub(k, big.NewInt(int64(wnaf[i])))
			length = i
		} else {
			wnaf[i] = 0
		}
		k.Rsh(k, 1)
	}
	length++
	wnafShort := make([]int8, length)
	copy(wnafShort, wnaf[0:length])
	return wnafShort
}

// Neg ---
func (e *ECPoint) Neg(x *ECPoint) *ECPoint {
	negY := NewECFieldElement()
	negY.Neg(x.Y)
	return &ECPoint{x.X, negY, x.curve}
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
	//vX := p.X.value.Bytes()
	//vY := p.Y.value.Bytes()
	//PrintHex("p point", []byte(vX), len(vX))
	//PrintHex("p point", []byte(vY), len(vY))

	k := big.NewInt(0)
	l := len(n)
	Tmp1 := make([]byte, l)

	copy(Tmp1, n)
	Reverse(Tmp1)
	//Tmp2 := make([]byte, l+1)
	//copy(Tmp2, Tmp1)

	k.SetBytes(Tmp1)
	kbts := k.Bytes()
	PrintHex("kbts", kbts, len(kbts))

	if k.Sign() == 0 {
		return Infinity
	}

	nP := Multiply(p, k)
	DumpECPoint(e, nP)
	return nP
}

// Add ---
func (e *ECPoint) Add(x *ECPoint, y *ECPoint) *ECPoint {
	if x.IsInfinity() {
		DumpECPoint(e, y)
		return y
	}
	if y.IsInfinity() {
		DumpECPoint(e, y)
		return x
	}
	if x.X.Equals(y.X) {
		if x.Y.Equals(y.Y) {
			twcX := x.Twice()
			DumpECPoint(e, twcX)
			return twcX
		}
		return Infinity
	}
	Tmp1 := NewECFieldElement()
	Tmp1.Sub(y.Y, x.Y)
	Tmp2 := NewECFieldElement()
	Tmp2.Sub(y.X, x.X)
	gama := NewECFieldElement()
	gama.Div(Tmp1, Tmp2)

	x3 := gama.Square()
	x3.Sub(x3, x.X)
	x3.Sub(x3, y.X)

	y3 := NewECFieldElement()
	y3.Sub(x.X, x3)
	y3.Mul(y3, gama)
	y3.Sub(y3, x.Y)

	ret := &ECPoint{x3, y3, x.curve}
	DumpECPoint(e, ret)
	return ret

}

// Sub ---
func (e *ECPoint) Sub(x *ECPoint, y *ECPoint) *ECPoint {
	if y.IsInfinity() {
		return x
	}
	tmp := NewECPoint()
	tmp.Neg(y)
	tmp.Add(tmp, x)

	DumpECPoint(e, tmp)

	return tmp
}
