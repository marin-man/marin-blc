package test

import (
	"bkc/core"
	"fmt"
	"testing"
)

func TestWallets_CreateWallet(t *testing.T) {
	wallets := core.NewWallets()
	wallets.CreateWallet()
	fmt.Printf("wallets:%v\n", wallets.Wallets)
}