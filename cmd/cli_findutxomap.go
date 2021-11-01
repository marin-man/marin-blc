package cmd

import (
	"bkc/core"
)

// TestResetUTXO 重置 utxo table
func (cli *CLI) TestResetUTXO(nodeId string) {
	blockchain := core.BlockchainObject(nodeId)
	defer blockchain.DB.Close()
	utxoSet := core.UTXOSet{Blockchain: blockchain}
	utxoSet.ResetUTXOSet()
}

// TestFindUTXOMap 查找
func (cli *CLI) TestFindUTXOMap()  {

}