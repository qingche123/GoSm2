package sm2

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

const (
	HASHLEN       = 32
	PRIVATEKEYLEN = 32
	PUBLICKEYLEN  = 32
	SIGNATURELEN  = 64
)

//-----------------------------------------------------------------------------

type ecdsaSignature struct {
	R, S *big.Int
}

//-----------------------------------------------------------------------------

func Init() {
	InitSecp256r1()
}

func Sha256(value []byte) []byte {
	//TODO: implement Sha256

	return nil
}

func RIPEMD160(value []byte) []byte {
	//TODO: implement RIPEMD160

	return nil
}

// Generate the "real" random number which can be used for crypto algorithm
func RandomNum(n int) ([]byte, error) {
	// TODO Get the random number from System urandom
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func Hash(data []byte) [HASHLEN]byte {
	return sha256.Sum256(data)
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func CaculateE(n *big.Int, message []byte) *big.Int {
	msgbitlen := len(message) * 8

	msgtmp := make([]byte, len(message)+1)
	copy(msgtmp, message)

	ReverseLen(msgtmp, len(message))

	trunc := big.NewInt(0).SetBytes(msgtmp)

	if n.BitLen() < msgbitlen {

		trunc.Rsh(trunc, uint(msgbitlen-n.BitLen()))
	}

	return trunc
}

//FIXME, does the privkey need base58 encoding?
//This generates a public & private key pair
func NewGenKeyPair() (*PrivateKey, *PublicKey, error) {

	mpubkey := new(PublicKey)
	mprikey := new(PrivateKey)

	k, _ := RandomNum(32)
	dk := big.NewInt(0).SetBytes(k)

	mprikey.d = dk
	mpubkey.pbkey.Mul(G, k)
	mprikey.pbkey = mpubkey.pbkey

	return mprikey, mpubkey, nil
}

func Sign(privateKey *PrivateKey, data []byte) ([]big.Int, error) {
	if privateKey == nil {
		fmt.Println("prikey is nil")
	}
	e := CaculateE(Ecurve.N, data)

	mpriKey := make([]byte, len(privateKey.d.Bytes())+1)
	copy(mpriKey, privateKey.d.Bytes())
	ReverseLen(mpriKey, len(privateKey.d.Bytes()))

	d := big.NewInt(0).SetBytes(mpriKey)

	r := big.NewInt(0)
	s := big.NewInt(0)

	for {
		k := big.NewInt(0)

		for {
			for {
				k, _ := rand.Prime(rand.Reader, Ecurve.N.BitLen())

				if k.Sign() == 0 || k.Cmp(Ecurve.N) >= 0 {
					break
				}
			}

			p := NewECPoint()

			p.Mul(G, k.Bytes())
			r.Mod(p.X.value, Ecurve.N)

			if r.Sign() != 0 {
				break
			}
		}

		Tmp := big.NewInt(0)
		Tmp1 := big.NewInt(0)
		Tmp2 := big.NewInt(0)

		Tmp.ModInverse(k, Ecurve.N)

		Tmp1.Mul(d, r)
		Tmp1.Add(Tmp1, e)

		Tmp.Mul(Tmp, Tmp1)
		Tmp.Mod(Tmp, Ecurve.N)

		Tmp2.Div(Ecurve.N, big.NewInt(2))

		if s.Cmp(Tmp2) == 1 {
			s.Sub(Ecurve.N, s)
		}

		if s.Sign() != 0 {
			break
		}
	}
	zz := make([]big.Int, 2)

	zz[0] = *r
	zz[1] = *s

	return zz, nil
}

func SumOfTwoMultiplies(P *ECPoint, k *big.Int, Q *ECPoint, l *big.Int) *ECPoint {
	m := 0
	if k.BitLen() > l.BitLen() {
		m = k.BitLen()
	} else {
		m = l.BitLen()
	}

	Z := NewECPoint()
	Z.Add(P, Q)

	R := NewECPoint()
	DumpECPoint(R, Infinity)

	for i := m - 1; i >= 0; i-- {
		R = R.Twice()

		if k.Bit(int(i)) == 1 {
			if l.Bit(int(i)) == 1 {
				R.Add(R, Z)
			} else {
				R.Add(R, P)
			}
		} else {
			if l.Bit(int(i)) == 1 {
				R.Add(R, Q)
			}
		}
	}
	return R
}

func Verify(message []byte, publicKey *PublicKey, r *big.Int, s *big.Int) bool {
	if r.Sign() < 1 || s.Sign() < 1 || r.Cmp(Ecurve.N) >= 0 || s.Cmp(Ecurve.N) >= 0 {
		return false
	}
	c := big.NewInt(0)
	u1 := big.NewInt(0)
	u2 := big.NewInt(0)

	e := CaculateE(Ecurve.N, message)
	c.ModInverse(s, Ecurve.N)

	u1.Mul(e, c)
	u1.Mod(u1, Ecurve.N)

	u2.Mul(r, c)
	u2.Mod(u1, Ecurve.N)

	point := SumOfTwoMultiplies(G, u1, publicKey.pbkey, u2)

	v := big.NewInt(0)
	v.Mod(point.X.value, Ecurve.N)

	return (v.Cmp(r) == 0)
}
