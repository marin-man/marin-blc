package cmd

import (
	"bkc/core"
	"fmt"
)

// createWallets 创建钱包集合
func (cli *CLI) createWallets(nodeId string) {
	wallets := core.NewWallets(nodeId) // 创建一个集合对象
	wallets.CreateWallet(nodeId)
	fmt.Println("当前的钱包信息")
	for key, _ := range wallets.Wallets {
		fmt.Printf("\t[%s]\n", key)
	}
}