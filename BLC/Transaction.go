package BLC

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"time"
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

// HashTransaction 生成交易哈希（交易序列化），不同时间生成的交易哈希值不同
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// 设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx Hash encoded failed %v\n", err)
	}
	// 添加时间戳标识，不添加会导致所有的 coinbase 交易哈希完全相同
	tm := time.Now().UnixNano()
	// 用于生成哈希的原数据
	txHashBytes := bytes.Join([][]byte{result.Bytes(), IntToHex(tm)}, []byte{})
	// 生成哈希值
	hash := sha256.Sum256(txHashBytes)
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

// Verity 验证签名
func (tx *Transaction) Verity(prevTxs map[string]Transaction) bool {
	// 检查能否找到交易哈希
	for _, vin := range tx.Vins {
		if prevTxs[hex.EncodeToString(vin.TxHash)].TxHash == nil {
			log.Panicf("VERIFY ERROR : transaction verity failed!\n")
		}
	}
	// 提取相同的交易签名属性
	txCopy := tx.TrimmedCopy()
	// 使用相同的椭圆
	curve := elliptic.P256()
	// 遍历 tx 输入，对每笔输入所引用的输出进行校验
	for vinId, vin := range tx.Vins {
		// 获取关联交易
		prevTx := prevTxs[hex.EncodeToString(vin.TxHash)]
		// 找到发送者（当前输入引用的哈希——输出的哈希）
		txCopy.Vins[vinId].PublicKey = prevTx.Vouts[vin.Vout].Ripemd160Hash
		// 由需要验证的数据生成的交易哈希，必须要与签名时的数据完全一致
		txCopy.TxHash = txCopy.Hash()
		// 在比特币中，签名是一个数值对，r、s 代表签名
		// 获取 r、s，两者长度相等，要从输入的 signature 中获取
		r, s := big.Int{}, big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen/2)])
		s.SetBytes(vin.Signature[(sigLen/2):])
		// 获取公钥，由 x，y 坐标组成
		x, y := big.Int{}, big.Int{}
		pubKeyLen := len(vin.PublicKey)
		x.SetBytes(vin.PublicKey[:(pubKeyLen/2)])
		y.SetBytes(vin.PublicKey[(pubKeyLen/2):])
		rawPublicKey := ecdsa.PublicKey{curve, &x, &y}
		// 调用验证签名核心函数
		if !ecdsa.Verify(&rawPublicKey, txCopy.TxHash, &r, &s) {
			return false
		}

	}
	return true
}