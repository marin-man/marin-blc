package BLC

import "bytes"

// 交易输入管理

// TxInput 输入结构
type TxInput struct {
	TxHash		[]byte		// 交易哈希（不是指当前的交易哈希）
	Vout		int			// 引用的上一笔交易的输出索引号
	Signature	[]byte		// 数字签名
	PublicKey	[]byte		// 公钥
}

// UnLockRipemd160Hash 传递哈希160进行判断
func (in *TxInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	// 获取 input 的ripemd160
	inputRipemd160Hash := Ripemd160Hash(in.PublicKey)
	return bytes.Compare(inputRipemd160Hash, ripemd160Hash) == 0
}