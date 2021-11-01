package core

import (
	"bkc/utils"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// 存入所有输出的集合

type TXOutputs struct {
	TXOutputs []*TxOutput
}

// TxOutput 交易的输出管理
type TxOutput struct {
	Value			int			// 金额
	Ripemd160Hash 	[]byte		//用脚本语言意味着比特币可以也作为智能合约平台
}

// UnLockScriptPubkeyWithAddress output 身份验证
func (out *TxOutput) UnLockScriptPubkeyWithAddress(address string) bool {
	// 转换
	hash160 := StringToHash160(address)
	fmt.Printf("%x\n", hash160)
	return bytes.Compare(hash160, out.Ripemd160Hash) == 0
}

// StringToHash160 string 转 hash160
func StringToHash160(address string) []byte {
	pubKeyHash := utils.Base58Decode([]byte(address))
	hash160 := pubKeyHash[:len(pubKeyHash) - addressCheckSumLen]
	return hash160[:]
}

// NewTxOutput 新建 output 对象
func NewTxOutput(value int, address string) *TxOutput {
	txOutput := &TxOutput{}
	hash160 := StringToHash160(address)
	txOutput.Value = value
	txOutput.Ripemd160Hash = hash160
	return txOutput
}

// Serialize 输出集合序列化
func (txOutputs *TXOutputs) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(txOutputs); nil != err {
		log.Panicf("serialize the utxo failed! %v\n", err)
	}
	return result.Bytes()
}

// Deserializer 输出集合序列化
func Deserializer(txOutputsBytes []byte) *TXOutputs {
	var txOutputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(txOutputsBytes))
	if err := decoder.Decode(&txOutputs); nil != err {
		log.Panicf("deserialize the struct utxo failed! %v\n", err)
	}
	return &txOutputs
}