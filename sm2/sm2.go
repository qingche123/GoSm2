package sm2

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// HASHLEN ---
const (
	HASHLEN       = 32
	PRIVATEKEYLEN = 32
	PUBLICKEYLEN  = 32
	SIGNATURELEN  = 64
)

type CurveParams struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	A       *big.Int // the constant of the curve equation
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

var BaseG *ECPoint
var Infinity *ECPoint
var EccParams *CurveParams

// Init ---
func Init() {
	EccParams = new(CurveParams)
	EccParams.P, _ = new(big.Int).SetString("8542D69E4C044F18E8B92435BF6FF7DE457283915C45517D722EDB8B08F1DFC3", 16)
	EccParams.N, _ = new(big.Int).SetString("8542D69E4C044F18E8B92435BF6FF7DD297720630485628D5AE74EE7C32E79B7", 16)
	EccParams.A, _ = new(big.Int).SetString("787968B4FA32C3FD2417842E73BBFEFF2F3C848B6831D7E0EC65228B3937E498", 16)
	EccParams.B, _ = new(big.Int).SetString("63E4C6D3B23B0C849CF84241484BFE48F61D59A5B16BA06E6E12D1DA27C5249A", 16)

	EccParams.Gx, _ = new(big.Int).SetString("421DEBD61B62EAB6746434EBC3CC315E32220B3BADD50BDC4C4E6C147FEDD43D", 16)
	EccParams.Gy, _ = new(big.Int).SetString("0680512BCBB42C07D47349D2153B70C4E5D7FDFCBFA36EA1A85841B9E46E09A2", 16)
	EccParams.BitSize = 256
	EccParams.Name = "sm2"

	Infinity = &ECPoint{nil, nil, EccParams}

	GX := &ECFieldElement{EccParams.Gx, EccParams}
	GY := &ECFieldElement{EccParams.Gy, EccParams}
	BaseG = &ECPoint{GX, GY, EccParams}

	return
}

// RandomNum Generate the "real" random number which can be used for crypto algorithm
func RandomNum(n int) ([]byte, error) {
	// TODO Get the random number from System urandom
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenKeyPair generates a public & private key pair
func GenKeyPair() ([]byte, *big.Int, *big.Int, error) {
	pubKey := NewECPoint()

	dBytes, _ := RandomNum(PRIVATEKEYLEN)
	pubKey.Mul(BaseG, dBytes)

	//PrintHex("prikey", dBytes, len(dBytes))
	//PrintHex("pbkeyX", pubKey.X.value.Bytes(), len(pubKey.X.value.Bytes()))
	//PrintHex("pbkeyY", pubKey.Y.value.Bytes(), len(pubKey.Y.value.Bytes()))

	return dBytes, pubKey.X.value, pubKey.Y.value, nil
}

// CaculateE ---
func CaculateE(curveN *big.Int, msg []byte) *big.Int {
	msgBitLen := len(msg) * 8

	trunc := new(big.Int).SetBytes(msg)

	if curveN.BitLen() < msgBitLen {
		trunc.Rsh(trunc, uint(msgBitLen-curveN.BitLen()))
	}
	return trunc
}

// Sign process:
// 1. choose an integer num k between 1 and n - 1.
// 2. compute point = k * BaseG.
// 3. compute r = (e + point.X) mod n, if r or r + k is equal 0 goto step 1.
// 4. compute ((1 + d)(-1) * (k - r*d))mod n, (-1) express modinverse operation
//    if s is equal 0 goto step 1.
//    e is the message, d is private key.
func Sign(priKey []byte, data []byte) (*big.Int, *big.Int, error) {
	if nil == priKey {
		fmt.Println("prikey is nil")
	}

	e := big.NewInt(0)
	e.SetBytes(data)

	priD := new(big.Int).SetBytes(priKey)

	k := big.NewInt(0)
	r := big.NewInt(0)
	s := big.NewInt(0)
	rAddK := big.NewInt(0)
	for {
		for {
			for {
				randK := make([]byte, EccParams.BitSize/8)
				_, err := rand.Read(randK)
				if err != nil {
					return nil, nil, err
				}
				//PrintHex("ranK", randK, len(randK))
				k.SetBytes(randK)
				if k.Sign() != 0 && k.Cmp(EccParams.N) < 0 {
					break
				}
			}

			kG := NewECPoint()
			kG.Mul(BaseG, k.Bytes())
			r.Add(e, kG.X.value)
			r.Mod(r, EccParams.N)

			if r.Sign() != 0 {
				rAddK.Add(r, k)
				if rAddK.Sign() != 0 {
					break
				}
			}
		}
		//s = ((1 + dA)-1 * (k - r*dA))mod n
		tmp := big.NewInt(0)
		tmp.Add(priD, big.NewInt(1))
		tmp.ModInverse(tmp, EccParams.N)

		tmp1 := big.NewInt(0)
		tmp1.Mul(r, priD)
		tmp1.Sub(k, tmp1)
		tmp1.Mod(tmp1, EccParams.N)

		s.Mul(tmp, tmp1)
		s.Mod(s, EccParams.N)

		if s.Sign() != 0 {
			break
		}
	}
	retR := big.NewInt(0)
	retS := big.NewInt(0)

	// r and s must between 1 and N - 1
	if r.Sign() < 1 {
		retR.Add(EccParams.P, r)
	} else {
		retR.Set(r)
	}

	if s.Sign() < 1 {
		retS.Add(EccParams.P, s)
	} else {
		retS.Set(s)
	}

	//PrintHex("r", sig[0].Bytes(), len(sig[0].Bytes()))
	//PrintHex("s", sig[1].Bytes(), len(sig[1].Bytes()))
	return retR, retS, nil
}

//Verify process:
// 1. computer t = (r' + s')mod n, if t = 0, verfy failed
// 2. computer (x1, y1) = [s']BaseG + [t]PubKey
// 3. computer R = (e +x1)
// 4. check that if R mod n == r, otherwise verify failed.
func Verify(X *big.Int, Y *big.Int, data []byte, r, s *big.Int) (bool, error) {
	if r.Sign() < 1 || s.Sign() < 1 || r.Cmp(EccParams.N) >= 0 || s.Cmp(EccParams.N) >= 0 {
		return false, errors.New("signature is invalid")
	}

	t := big.NewInt(0)
	t.Add(r, s)
	t.Mod(t, EccParams.N)

	pub := NewECPoint()
	pub.X.value.Set(X)
	pub.Y.value.Set(Y)

	point := SumOfTwoMultiplies(BaseG, s, pub, t)

	e := new(big.Int).SetBytes(data)
	R := big.NewInt(0)
	R.Add(e, point.X.value)
	R.Mod(R, EccParams.N)

	return (0 == R.Cmp(r)), nil
}
