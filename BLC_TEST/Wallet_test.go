package BLC_TEST

import (
	"bkc/BLC"
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := BLC.NewWallet()
	fmt.Printf("private key : %v\n", wallet.PrivateKey)
	fmt.Printf("public key : %v\n", wallet.PublicKey)
	fmt.Printf("wallet : %v\n", wallet)
}

func TestWallet_GetAddress(t *testing.T) {
	wallet := BLC.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("the address of coin is [%s]\n", address)
	fmt.Printf("the validation of current address is %v\n", BLC.IsValidForAddress([]byte(address)))
}
