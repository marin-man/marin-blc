package test

import (
	"bkc/core"
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := core.NewWallet()
	fmt.Printf("private key : %v\n", wallet.PrivateKey)
	fmt.Printf("public key : %v\n", wallet.PublicKey)
	fmt.Printf("wallet : %v\n", wallet)
}

func TestWallet_GetAddress(t *testing.T) {
	wallet := core.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("the address of coin is [%s]\n", address)
	fmt.Printf("the validation of current address is %v\n", core.IsValidForAddress([]byte(address)))
}
