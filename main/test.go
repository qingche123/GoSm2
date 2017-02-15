package main

import (
	. "GoSm2/sm2"
	"fmt"
	"os"
)

func main() {
	Init()

	prikey, pubkey, err := NewGenKeyPair()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buf := "This is a message to be signed and verified by ECDSA!"
	// Sign ecdsa style

	signature, err := Sign(prikey, []byte(buf))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Signature : %x\n", signature)

	// Verify
	verifystatus := Verify([]byte(buf), pubkey, &signature[0], &signature[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(verifystatus) // should be true
}
