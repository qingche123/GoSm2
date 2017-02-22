package main

import (
	"GoSm2/sm2"
	"fmt"
	"math/big"
	"os"
)

func PrintHex(str string, bt []byte, length int) {
	fmt.Println(str, "Length = ", length)
	for i := 0; i < length; i++ {
		if i%16 == 0 && i != 0 {
			fmt.Println()
		}
		fmt.Printf("0x%02x,", bt[i])
	}
	fmt.Println(" ")
	fmt.Println(" ")
}

func reverse(data []byte) {
	len1 := len(data)
	for i := 0; i < len1/2; i++ {
		tmp := data[i]
		data[i] = data[len1-i-1]
		data[len1-i-1] = tmp
	}
}

func test() {

	z, _ := new(big.Int).SetString("917049492264781353612408419258525201475102643932529502393738587225651523777317", 10)
	sz := z.Bytes()

	dBytes := make([]byte, 32)
	copy(dBytes, sz)
	//reverse(dBytes)

	PrintHex("dBytes", dBytes, len(dBytes))
	/*
		a := big.NewInt(1*16 + 2*16*16 + 3*16*16*16)
		bt := a.Bytes()
		PrintHex("a", bt, len(bt))

		v := []byte{0x32, 0x10}
		b := big.NewInt(0)
		b.SetBytes(v)
		fmt.Println(b.Int64())

		sv := b.Bytes()
		PrintHex("sv", sv, len(sv))*/
}

func main() {
	//test()
	//sm2.TestPointEncode()

	sm2.Init()

	prikey, pubkey, err := sm2.NewGenKeyPair()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := "This is a message to be signed and verified by ECDSA!"
	// Sign ecdsa style

	signature, err := sm2.Sign(prikey, []byte(buf))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Signature : %x\n", signature)

	// Verify
	verifystatus := sm2.Verify([]byte(buf), pubkey, &signature[0], &signature[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(verifystatus) // should be true
}
