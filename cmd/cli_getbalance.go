package cmd

import (
	"bkc/core"
	"fmt"
)

// getBalance 查询余额
func (cli *CLI) getBalance(from string, nodeId string) {
	// 查找该地址 UTXO
	// 获取区块链对象
	blockchain := core.BlockchainObject(nodeId)
	defer blockchain.DB.Close()   // 关闭实例对象
	utxoSet := core.UTXOSet{Blockchain: blockchain}
	amount := utxoSet.GetBalance(from)
	fmt.Printf("地址 [%s] 的余额：[%d]\n", from, amount)
}