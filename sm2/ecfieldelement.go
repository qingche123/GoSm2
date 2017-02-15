package sm2

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type ECFieldElement struct {
	value *big.Int
	curve *ECCurveParams
}

func NewECFieldElement() *ECFieldElement {

	bnx := big.NewInt(0)

	return &ECFieldElement{bnx, Ecurve}
}

func DumpECFieldElement(dst *ECFieldElement, src *ECFieldElement) {

	dst.value.Set(src.value)
	dst.curve = src.curve
}

func GetLowestSetBit(k *big.Int) int {
	i := 0
	for i = 0; k.Bit(i) != 1; i++ {
	}
	return i
}

func FastLucasSequence(A, B, C, k *big.Int) []big.Int {
	n := k.BitLen()
	s := GetLowestSetBit(k)

	Uh := big.NewInt(1)
	Vl := big.NewInt(2)
	Ql := big.NewInt(1)
	Qh := big.NewInt(1)
	Vh := big.NewInt(0)
	Tmp := big.NewInt(0)

	Vh.Set(B)

	for j := n - 1; j >= s+1; j-- {
		Ql.Mul(Ql, Qh)
		Ql.Mod(Ql, A)

		if k.Bit(j) == 1 {
			Qh.Mul(Ql, C)
			Qh.Mod(Qh, A)

			Uh.Mul(Uh, Vh)
			Uh.Mod(Uh, A)

			Vl.Mul(Vh, Vl)
			Tmp.Mul(B, Ql)
			Vl.Sub(Vl, Tmp)
			Vl.Mod(Vl, A)

			Vh.Mul(Vh, Vh)
			Tmp.Lsh(Qh, 1)
			Vh.Sub(Vh, Tmp)
			Vh.Mod(Vh, A)

		} else {
			Qh.Set(Ql)

			Uh.Mul(Uh, Vl)
			Uh.Sub(Uh, Ql)
			Uh.Mod(Uh, A)

			Vh.Mul(Vh, Vl)
			Tmp.Mul(B, Ql)
			Vh.Sub(Vh, Tmp)
			Vh.Mod(Vh, A)

			Vl.Mul(Vl, Vl)
			Tmp.Lsh(Ql, 1)
			Vl.Sub(Vl, Tmp)
			Vl.Mod(Vl, A)
		}
	}

	Ql.Mul(Ql, Qh)
	Ql.Mod(Ql, A)

	Qh.Mul(Ql, C)
	Qh.Mod(Qh, A)

	Uh.Mul(Uh, Vl)
	Uh.Sub(Uh, Ql)
	Uh.Mod(Uh, A)

	Vl.Mul(Vh, Vl)
	Tmp.Mul(B, Ql)
	Vl.Sub(Vl, Tmp)
	Vl.Mod(Vl, A)

	Ql.Mul(Ql, Qh)
	Ql.Mod(Ql, A)

	for j := 1; j <= s; j++ {
		Uh.Mul(Uh, Vl)
		Uh.Mul(Uh, A)

		Vl.Mul(Vl, Vl)
		Tmp.Lsh(Ql, 1)
		Vl.Sub(Vl, Tmp)
		Vl.Mod(Vl, A)

		Ql.Mul(Ql, Ql)
		Ql.Mod(Ql, A)
	}

	//var Array [2]*big.Int
	//Array = [2]*big.Int{Uh, Vl}

	bnret := make([]big.Int, 2)
	bnret[0] = *Uh
	bnret[1] = *Vl

	return bnret
}

func IsEven(k *big.Int) bool {
	z := big.NewInt(0)
	z.Mod(k, big.NewInt(2))
	if z.Int64() == 0 {
		return true
	} else {
		return false
	}
}

func Reverse(data []byte) {

	len1 := len(data)

	for i := 0; i < len1/2; i++ {
		Tmp := data[i]
		data[i] = data[len1-1-i]
		data[len1-1-i] = Tmp
	}
}

func ReverseLen(data []byte, length int) {

	for i := 0; i < length/2; i++ {
		Tmp := data[i]
		data[i] = data[length-1-i]
		data[length-1-i] = Tmp
	}
}

func (e *ECFieldElement) CompareTo(other *ECFieldElement) int {
	if e == other {
		return 0
	}
	return e.value.Cmp(other.value)
}

func (e *ECFieldElement) Equals(other *ECFieldElement) bool {
	if e == other {
		return true
	}
	if other == nil {
		return false
	}
	return (e.value.Cmp(other.value) == 0)
}

func (e *ECFieldElement) Square() *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Mul(e.value, e.value)
	Tmp.Mod(Tmp, e.curve.P)

	return &ECFieldElement{Tmp, e.curve}
}

func (e *ECFieldElement) Sqrt() *ECFieldElement {
	if e.curve.P.Bit(1) == 1 {
		Tmp1 := big.NewInt(0)
		Tmp1.Rsh(e.curve.P, 2)
		Tmp1.Add(Tmp1, big.NewInt(1))

		Tmp2 := big.NewInt(0)
		Tmp2.Exp(e.value, Tmp1, e.curve.P)

		z := &ECFieldElement{Tmp2, e.curve}

		if z.Square().Equals(e) {
			return z
		} else {
			fmt.Println("error z^2 != z")
			return nil
		}
	}

	qMinusOne := big.NewInt(0)
	qMinusOne.Sub(e.curve.P, big.NewInt(1))

	legendExponent := big.NewInt(0)
	legendExponent.Rsh(qMinusOne, 1)

	Tmp := big.NewInt(0)
	Tmp.Exp(e.value, legendExponent, e.curve.P)
	if Tmp.Cmp(big.NewInt(1)) != 0 {
		return nil
	}

	u := big.NewInt(0)
	u.Rsh(qMinusOne, 2)

	k := big.NewInt(0)
	k.Lsh(u, 1)
	k.Add(k, big.NewInt(1))

	Q := big.NewInt(0)
	Q.Set(e.value)
	fourQ := big.NewInt(0)
	fourQ.Lsh(Q, 2)
	fourQ.Mod(fourQ, e.curve.P)

	U := big.NewInt(0)
	V := big.NewInt(0)

	for {
		P := big.NewInt(0)
		for {
			Tmp1 := big.NewInt(0)
			P, _ := rand.Prime(rand.Reader, e.curve.P.BitLen())

			if P.Cmp(e.curve.P) < 0 {
				Tmp1.Mul(P, P)
				Tmp1.Sub(Tmp1, fourQ)
				Tmp1.Exp(Tmp1, legendExponent, e.curve.P)

				if Tmp1.Cmp(qMinusOne) == 0 {
					break
				}
			}
		}

		//var Array [2]*big.Int
		result := FastLucasSequence(e.curve.P, P, Q, k)

		U.Set(&result[0])
		V.Set(&result[1])

		Tmp2 := big.NewInt(0)
		Tmp2.Mul(V, V)
		Tmp2.Mod(Tmp2, e.curve.P)
		if Tmp2.Cmp(fourQ) == 0 {
			if V.Bit(0) == 1 {
				V.Add(V, e.curve.P)
			}
			V.Rsh(V, 1)

			return &ECFieldElement{V, e.curve}
		}
		if (U.Cmp(big.NewInt(0)) != 0) || (U.Cmp(qMinusOne) != 0) {
			break
		}
	}
	return nil
}

func (e *ECFieldElement) ToByteArray() []byte {

	data := e.value.Bytes()
	dlen := len(data)

	if len(data) == 32 {
		Reverse(data)
		return data
	}
	if len(data) > 32 {
		data1 := data[0:32]
		Reverse(data1)
		return data1
	}
	Reverse(data)
	var data2 = make([]byte, 32)

	copy(data2[32-dlen:], data)
	return data2
}

func (e *ECFieldElement) Neg(x *ECFieldElement) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Neg(e.value)
	Tmp.Mod(Tmp, x.curve.P)

	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) Mul(x *ECFieldElement, y *ECFieldElement) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Mul(x.value, y.value)
	Tmp.Mod(Tmp, x.curve.P)
	e.value.Set(Tmp)

	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) MulBig(x *ECFieldElement, y *big.Int) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Mul(x.value, y)
	Tmp.Mod(Tmp, x.curve.P)
	e.value.Set(Tmp)

	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) Div(x *ECFieldElement, y *ECFieldElement) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp1 := big.NewInt(0)
	Tmp1.ModInverse(y.value, x.curve.P)
	Tmp.Mul(x.value, Tmp1)
	Tmp.Mod(Tmp, x.curve.P)

	e.value.Set(Tmp)
	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) DivBig(x *ECFieldElement, y *big.Int) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp1 := big.NewInt(0)
	Tmp1.ModInverse(y, x.curve.P)
	Tmp.Mul(x.value, Tmp1)
	Tmp.Mod(Tmp, x.curve.P)

	e.value.Set(Tmp)
	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) Add(x *ECFieldElement, y *ECFieldElement) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Add(x.value, y.value)
	Tmp.Mod(Tmp, x.curve.P)
	e.value.Set(Tmp)
	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) AddBig(x *ECFieldElement, y *big.Int) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Add(x.value, y)
	Tmp.Mod(Tmp, x.curve.P)
	e.value.Set(Tmp)
	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) Sub(x *ECFieldElement, y *ECFieldElement) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Sub(x.value, y.value)
	Tmp.Mod(Tmp, x.curve.P)
	e.value.Set(Tmp)
	return &ECFieldElement{Tmp, x.curve}
}

func (e *ECFieldElement) SubBig(x *ECFieldElement, y *big.Int) *ECFieldElement {

	Tmp := big.NewInt(0)
	Tmp.Sub(x.value, y)
	Tmp.Mod(Tmp, x.curve.P)
	e.value.Set(Tmp)
	return &ECFieldElement{Tmp, x.curve}
}
