package core

// UTXO 结构管理
type UTXO struct {
	// UTXO 对应的交易哈希
	TxHash	[]byte
	// UTXO 在其所属交易的输出列表中的索引
	Index	int
	// Output 本身
	Output	*TxOutput
}
