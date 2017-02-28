package main

import (
	"GoSm2/sm2"
	"fmt"
	"os"
	"time"
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

func main() {
	sm2.Init()
	buf := "This is a message to be signed and verified by SM2!"

	priKey, x, y, err := sm2.GenKeyPair()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t1 := time.Now()
	for i := 0; i < 10; i++ {
		r, s, serr := sm2.Sign(priKey, []byte(buf))
		if nil != serr {
			fmt.Println(serr)
			os.Exit(1)
		}

		PrintHex("R", r.Bytes(), len(r.Bytes()))
		PrintHex("S", s.Bytes(), len(s.Bytes()))

		status, _ := sm2.Verify(x, y, []byte(buf), r, s)
		if true != status {
			fmt.Println("Verify Failed")
			os.Exit(1)
		}
		fmt.Println(status) // should be true
	}

	t2 := time.Now()
	d := t2.Sub(t1)
	fmt.Println(d)

}
