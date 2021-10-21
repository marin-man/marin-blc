package BLC_TEST

import (
	"bkc/BLC"
	"fmt"
	"testing"
)

func TestWallets_CreateWallet(t *testing.T) {
	wallets := BLC.NewWallets()
	wallets.CreateWallet()
	fmt.Printf("wallets:%v\n", wallets.Wallets)
}