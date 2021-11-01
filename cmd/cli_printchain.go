package cmd

import (
	"bkc/core"
	"fmt"
	"os"
)

// printChain 打印完整区块链信息
func (cli *CLI) printChain(nodeId string) {
	if !core.DBExits(nodeId) {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := core.BlockchainObject(nodeId)
	blockchain.PrintChain()
}