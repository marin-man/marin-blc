package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
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
		Signature: nil,
		PublicKey: nil,
	}
	// 输出：value，address
	txoOutput := NewTxOutput(10, address)

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
	// 获取钱包集合对象
	wallets := NewWallets()
	wallet := wallets.Wallets[from]
	// 输入
	for txHash, indexArray := range spendableUTXODic {
		txHashBytes, err := hex.DecodeString(txHash)
		if nil != err {
			log.Panicf("decode string to []byte failed! %v\n", err)
		}
		// 遍历索引列表
		for _, index := range indexArray {
			txInput := &TxInput{txHashBytes, index, nil, wallet.PublicKey}
			txInputs = append(txInputs, txInput)
		}
	}
	// 输出（转账源）
	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)
	// 找零
	if money < amount {
 		txOutput = NewTxOutput(money - amount, from)
		txOutputs = append(txOutputs, txOutput)
	} else {
		log.Panicf("余额不足...\n")
	}

	tx := Transaction{
		TxHash: nil,
		Vins: txInputs,
		Vouts: txOutputs,
	}
	// 生成一笔完整的交易
	tx.HashTransaction()
	// 对交易进行签名
	bc.SignTransaction(&tx, wallet.PrivateKey)
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

// Sign 交易签名
// prevTxs：代表当前交易的输入所引用的所有 OUTPUT 所属的交易
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	// 处理输入，保证交易的正确性
	// 检查 tx 中每一个输入所引用的交易哈希是否包含在 prevTxs 中
	// 如果没有包含在里面，则说明该交易被人修改了
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("ERROR: Prev transaction is no correct!\n")
		}
	}
	// 提取需要签名的属性
	txCopy := tx.TrimmedCopy()
	// 处理交易副本的输入
	for vin_id, vin := range txCopy.Vins {
		// 获取关联交易
		prevTxs := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 找到发送者（当前输入引用的哈希——输出的哈希）
		txCopy.Vins[vin_id].PublicKey = prevTxs.Vouts[vin.Vout].Ripemd160Hash
		// 生成交易副本的哈希
		txCopy.TxHash = txCopy.Hash()
		// 调用核心签名函数
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxHash)
		if nil != err {
			log.Panicf("sign to transaction [%x] failed！ %v\n", txCopy, err)
		}
		// 组成交易签名
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vins[vin_id].Signature = signature
	}
}

// TrimmedCopy 交易拷贝，生成一个专门用于交易签名的副本
func (tx *Transaction) TrimmedCopy() Transaction {
	// 重新组装生成一个新的交易
	var inputs []*TxInput
	var outputs []*TxOutput
	// 组装 input
	for _, vin := range tx.Vins {
		inputs = append(inputs, &TxInput{
			vin.TxHash,
			vin.Vout,
			nil,
			nil,
		})
	}
	// 组装 output
	for _, vout := range tx.Vouts {
		outputs = append(outputs, &TxOutput{
			vout.Value,
			vout.Ripemd160Hash,
		})
	}
	txCopy := Transaction{tx.TxHash, inputs, outputs}
	return txCopy
}

// Serialize 交易序列化
func (tx *Transaction) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(tx); nil != err {
		log.Panicf("serialize the tx to []byte failed! %v\n", err)
	}
	return buffer.Bytes()
}

// Hash 设置用于签名的交易的哈希
func (tx *Transaction) Hash() []byte {
	txCopy := tx
	txCopy.TxHash = []byte{}
	hash := sha256.Sum256(tx.Serialize())
	return hash[:]
}