package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

// UTXO 持久化相关管理

// 用于存入 utxo 的 bucket
const utxoTableName = "utxoTable"

// UTXOSet 结构（保存指定区块中所有的 UTXO）
type UTXOSet struct {
	Blockchain	*BlockChain
}

// ResetUTXOSet 重置
func (utxoSet *UTXOSet) ResetUTXOSet() {
	// 在第一次创建时就更新 utxo table
	utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			err := tx.DeleteBucket([]byte(utxoTableName))
			if nil != err {
				log.Panicf("delete the utxo table failed! %v\n", err)
			}
		}

		// 创建
		bucket, err := tx.CreateBucket([]byte(utxoTableName))
		if nil != err {
			log.Panicf("create bucket failed! %v\n", err)
		}
		if nil != bucket {
			// 查找当前所有 UTXO
			txOutputMap := utxoSet.Blockchain.FindUTXOMap()

			for keyHash, outputs := range txOutputMap {
				// 将所有 UTXO 存入
				txHash, _ := hex.DecodeString(keyHash)
				fmt.Println("keyhash: %v\n", txHash)
				// 存入 utxo table
				err := bucket.Put(txHash, outputs.Serialize())
				if nil != err {
					log.Panicf("put the utxo into table failed! %v\n", err)
				}
			}
		}
		return nil
	})
}

// FindUTXOWithAddress 查找
func (utxoSet *UTXOSet) FindUTXOWithAddress(address string) []*UTXO {
	var utxos []*UTXO
	err := utxoSet.Blockchain.DB.View(func(tx *bolt.Tx) error {
		// 1. 获取 utxotable
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			// cursor 游标
			c := b.Cursor()
			// 通过游标白能力 boltdb 数据库中的数据
			for k, v := c.First(); nil != k; k, v = c.Next() {
				txOutputs := Deserializer(v)
				for _, utxo := range txOutputs.TXOutputs {
					if utxo.UnLockScriptPubkeyWithAddress(address) {
						utxo_signle := UTXO{Output: utxo}
						utxos = append(utxos, &utxo_signle)
					}
				}
			}
		}
		return nil
	})
	if nil != err {
		log.Panicf("find the utxo of [%s] failed! %v\n", address, err)
	}
	return utxos
}

// GetBalance 查询余额
func (utxoSet *UTXOSet) GetBalance(address string) int {
	UTXOS := utxoSet.FindUTXOWithAddress(address)
	var amount int
	for _, utxo := range UTXOS {
		fmt.Printf("utxo-txhash: %x\n", utxo.TxHash)
		fmt.Printf("utxo-index: %x\n", utxo.Index)
		fmt.Printf("utxo-Ripemd160Hash: %x\n", utxo.Output.Ripemd160Hash)
		fmt.Printf("utxo-Value: %x\n", utxo.Output.Value)
		amount += utxo.Output.Value
	}
	return amount
}


// 更新
func (utxoSet *UTXOSet) update() {
	// 获取最新区块
	latest_block := utxoSet.Blockchain.Iterator().Next()
	utxoSet.Blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoTableName))
		if nil != b {
			// 查找最新一个区块的交易列表，因为每上链一个区块，utxo table 都更新一次，所以只需要查找最近一个区块中的交易
			for _, tx := range latest_block.Txs {
				if !tx.IsCoinbaseTransaction() {
					// 2. 将已经被当前这笔交易的输入所引用的 UTXO 删掉
					for _, vin := range tx.Vins {
						// 需要更新的输出
						updateOutputs := TXOutputs{}
						// 获取指定输入所引用的交易哈希的输出
						outputBytes := b.Get(vin.TxHash)
						// 输出列表
						outs := Deserializer(outputBytes)
						for outIdx, out := range outs.TXOutputs {
							if vin.Vout != outIdx {
								updateOutputs.TXOutputs = append(updateOutputs.TXOutputs, out)
							}
						}
						// 如果交易中没有 UTXO 了，删除该交易
						if len(updateOutputs.TXOutputs) == 0 {
							b.Delete(vin.TxHash)
						} else {
							// 将更新之后的 utxo 数据存入数据库
							b.Put(vin.TxHash, updateOutputs.Serialize())
						}
					}

				}

				// 获取当前区块中新生成的交易输出
				// 1. 将最新区块中的 UTXO 插入
				newOutputs := TXOutputs{}
				newOutputs.TXOutputs = append(newOutputs.TXOutputs, tx.Vouts...)
				b.Put(tx.TxHash, newOutputs.Serialize())
			}
		}
		return nil
	})
}