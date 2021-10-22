package BLC

import "bytes"

// TxOutput 交易的输出管理
type TxOutput struct {
	Value			int			// 金额
	Ripemd160Hash	[]byte		// 用户名（UTXO 的所有者）
}

// UnLockScriptPubkeyWithAddress output 身份验证
func (ou *TxOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	// 转换
	hash160 := StringToHash160(address)
	return bytes.Compare(hash160, ou.Ripemd160Hash) == 0
}

// StringToHash160 string 转 hash160
func StringToHash160(address string) []byte {
	pubKeyHash := Base58Decode([]byte(address))
	hash160 := pubKeyHash[:len(pubKeyHash) - addressCheckSumLen]
	return hash160[:]
}

// 新建 output 对象
func NewTxOutput(value int, address string) *TxOutput {
	txOutput := &TxOutput{
		value,
		StringToHash160(address),
	}
	return txOutput
}