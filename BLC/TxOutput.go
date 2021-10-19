package BLC

// TxOutput 交易的输出管理
type TxOutput struct {
	Value			int			// 金额
	ScriptPubkey	string		// 用户名
}

// CheckPubkeyWithAddress 验证当前 UTXO 是否属于指定的地址
func (txOutput *TxOutput) CheckPubkeyWithAddress(address string) bool {
	return address == txOutput.ScriptPubkey
}