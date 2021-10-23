package BLC

import (
	"fmt"
	"os"
)

// send 发起交易
func (cli *CLI) send(from, to, amount []string)  {
	if !dbExit() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	// 获取区块链对象
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	if len(from) != len(to) && len(from) != len(amount) {
		fmt.Println("交易参数输入有误，请检查一致性...")
		os.Exit(1)
	}
	// 发起交易，生成新的区块
	blockchain.MineNewBlock(from, to, amount)
	// 调用 utxo table 的函数，更新 utxo table
	utxoSet := &UTXOSet{Blockchain: blockchain}
	utxoSet.update()
}