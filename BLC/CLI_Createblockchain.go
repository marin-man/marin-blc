package BLC

// createBlockchain 初始化区块链
func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockChain(address)
	defer bc.DB.Close()

	// 设置 utxo 重置操作
	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}