package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

// 交易管理文件

// Transaction 定义一个交易基本结构
type Transaction struct {
	TxHash		[]byte			// 交易哈希标识
	Vins		[]*TxInput		// 输入列表
	Vouts		[]*TxOutput		// 输出列表
}

// NewCoinbaseTransaction 实现 coinbase 交易
func NewCoinbaseTransaction(address string) *Transaction {
	var txCoinbase *Transaction
	// 输入，coinbase 特点：
	// txHash: nil, vout: -1, ScriptSig: 系统奖励
	txInput := &TxInput{
		TxHash: []byte{},
		Vout: -1,
		ScriptSig: "system reward",
	}
	// 输出：value，address
	txoOutput := &TxOutput{
		Value: 10,
		ScriptPubkey: address,
	}
	// 输入输出组装交易
	txCoinbase = &Transaction{
		TxHash: nil,
		Vins: []*TxInput{txInput},
		Vouts: []*TxOutput{txoOutput},
	}
	// 交易哈希生成
	txCoinbase.HashTransaction()
	return txCoinbase
}

// NewSimpleTransaction 生成普通转账交易
func NewSimpleTransaction(from string, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput		// 输入列表
	var txOutputs []*TxOutput	// 输出列表
	// 调用可花费 UTXO 函数
	money, spendableUTXODic := bc.FindSpendableUTXO(from, amount, txs)
	// 输入
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, err := hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string to []byte failed! %v\n", err)
		}
		// 遍历索引列表
		for _, index := range indexArray {
			txInput := &TxInput{index, from, txHashBytes}
			txInputs = append(txInputs, txInput)
		}
	}
	// 输出（转账源）
	txOutput := &TxOutput{}
	txOutputs = append(txOutputs, txOutput)
	// 找零
	if money < amount {
		txOutput = &TxOutput{money - amount, from}
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足...\n")
	}

	tx := Transaction{
		TxHash: nil,
		Vins: txInputs,
		Vouts: txOutputs,
	}
	tx.HashTransaction()
	return &tx
}

// HashTransaction 生成交易哈希（交易序列化）
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// 设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed %v\n", err)
	}

	// 生成哈希值
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

// IsCoinbaseTransaction 判断指定的交易是否时一个 coinbase 交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return -1 == tx.Vins[0].Vout && 0 == len(tx.Vins[0].TxHash)
}