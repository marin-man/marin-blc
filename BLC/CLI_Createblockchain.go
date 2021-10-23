package BLC

// createBlockchain 初始化区块链
func (cli *CLI) createBlockchain(address string, nodeId string) {
	bc := CreateBlockChain(address, nodeId)
	defer bc.DB.Close()

	// 设置 utxo 重置操作
	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}