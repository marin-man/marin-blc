package BLC

import (
	"fmt"
	"os"
)

// printChain 打印完整区块链信息
func (cli *CLI) printChain() {
	if !dbExit() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	blockchain.PrintChain()
}