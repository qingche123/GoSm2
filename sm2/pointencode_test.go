package sm2

import (
	"fmt"
	"testing"
)

func PrintHex(str string, bt []byte, length int) {
	fmt.Println(str, "Length = ", length)
	for i := 0; i < length; i++ {
		fmt.Printf("0x%02x ", bt[i])
	}
	fmt.Println(" ")
	fmt.Println(" ")
}
func TestPointEncode(t *testing.T) {

	Init()

	k1, _ := RandomNum(32)
	k2, _ := RandomNum(32)

	PrintHex("k1: ", k1, len(k1))
	PrintHex("k2: ", k2, len(k2))

	ap := NewECPoint()

	ap.X.value.SetBytes(k1)
	ap.Y.value.SetBytes(k2)

	fmt.Println("ap.X: ", ap.X.value, "\n")
	fmt.Println("ap.Y: ", ap.Y.value, "\n")

	ap.X.curveParam = Ecurve
	ap.Y.curveParam = Ecurve

	ap.curve = Ecurve

	encAp := ap.EncodePoint(true)
	PrintHex("encAp: ", encAp, len(encAp))

	bp := DecodePoint(encAp, Ecurve)
	fmt.Println("bp.X: ", bp.X.value, "\n")
	fmt.Println("bp.Y: ", bp.Y.value, "\n")
}
