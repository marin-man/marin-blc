package cmd

import (
	"bkc/core"
	"fmt"
)

// GetAccounts 获取钱包地址列表
func (cli *CLI) GetAccounts(nodeId string) {
	wallets := core.NewWallets(nodeId)
	fmt.Println("账号列表")
	for address := range wallets.Wallets {
		fmt.Printf("\t[%s]\n", address)
	}
}