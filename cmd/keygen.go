package main

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/nacl/box"
	b64 "encoding/base64"
)

func main() {

	publicKey1, privateKey1, _ := box.GenerateKey(rand.Reader)
	fmt.Println("publicKey1: ", myStringBase64(publicKey1))
	fmt.Println("privateKey1: ", myStringBase64(privateKey1))

}

func myStringBase64(bbb *[32]byte) string {
	var bar []byte = bbb[:]
	return b64.StdEncoding.EncodeToString(bar)
}
