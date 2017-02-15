package sm2

import (
	"math/big"
)

type PublicKey struct {
	pbkey *ECPoint
}

// PrivateKey represents a ECDSA private key.
type PrivateKey struct {
	pbkey *ECPoint
	d     *big.Int
}
