package BLC

// 交易输入管理

// TxInput 输入结构
type TxInput struct {
	Vout		int			// 引用的上一笔交易的输出索引号
	ScriptSig	string		// 用户名
	TxHash		[]byte		// 交易哈希（不是指当前的交易哈希）
}

// 验证引用的地址是否匹配
func (txInput *TxInput) CheckPubkeyWithAddress(address string) bool {
	return address == txInput.ScriptSig
}
