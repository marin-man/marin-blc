package BLC

import (
	"fmt"
	"os"
)

// printChain 打印完整区块链信息
func (cli *CLI) printChain(nodeId string) {
	if !dbExit(nodeId) {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject(nodeId)
	blockchain.PrintChain()
}