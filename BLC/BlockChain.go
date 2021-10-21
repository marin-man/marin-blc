package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strconv"
)

// 区块链管理工具
// 数据库名称
const dbName = "block.db"
// 表名称
const blockTableName = "blocks"

// BlockChain 区块链的基本结构
type BlockChain struct {
	DB		*bolt.DB	// 数据库对象
	Tip		[]byte		// 保存最新区块的哈希值
}

// CreateBlockChain 初始化区块链
func CreateBlockChain(address string) *BlockChain{
	// 保存最新区块的哈希值
	var blockHash []byte
	// 创建或打开一个数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("open db [%s] failed %v \n", dbName, err)
	}
	// 创建桶（表）,把创世区块存入数据库
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil == b {
			// 没找到桶
			b, err := tx.CreateBucket([]byte(blockTableName))
			if nil != err {
				log.Panicf("create bucket [%s] failed %v\n", blockTableName, err)
			}
			// 生成一个 coinbase 交易
			txCoinbase := NewCoinbaseTransaction(address)
			// 创建一个创世块
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
			// 存储
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if nil != err {
				log.Panicf("insert the genesis block failed %v\n", err)
			}
			blockHash = genesisBlock.Hash
			// 存钱最新区块的哈希
			err = b.Put([]byte("1"), genesisBlock.Hash)
			if nil != err {
				log.Panicf("saave the hash of genesis block failed %v\n", err)
			}
		} else {
			blockHash = b.Get([]byte("1"))
		}
		return nil
	})
	return &BlockChain{
		DB: db,
		Tip: blockHash,
	}
}

// AddBlock 添加区块到区块链中
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// 更新区块数据(insert)
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取数据库桶
		b := tx.Bucket([]byte(blockTableName))
		if nil != b{
			// 获取最后插入的区块
			blockBytes := b.Get(bc.Tip)
			// 区块数据反序列化
			latestBlock := Deserialize(blockBytes)
			// 新建区块
			newBlock := NewBlock(latestBlock.Height + 1, latestBlock.Hash, txs)
			// 存入数据库
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if nil != err {
				log.Panicf("insert new block to db failed %v\n", err)
			}
			// 更新最新区块的哈希
			err = b.Put([]byte("1"), newBlock.Hash)
			if nil != err {
				log.Panicf("update the latest block hash to db failed %v\n", err)
			}
			bc.Tip = newBlock.Hash
		}
		return nil
	})

	if nil != err {
		log.Panicf("insert block to db failed %v\n", err)
	}
}

// PrintChain 遍历数据库，输出所有区块信息
func (bc *BlockChain) PrintChain() {
	bcit := bc.Iterator()
	fmt.Println("区块链完整信息...")
	// 读取数据库
	// 循环读取
	for {
		if bcit.HasNext() {
			fmt.Println("-------------------------------")
			curBlock := bcit.Next()
			fmt.Printf("\tHash:%x\n", curBlock.Hash)
			fmt.Printf("\tPrevBlockHash:%x\n", curBlock.PrevBlockHash)
			fmt.Printf("\tTimeStamp:%v\n", curBlock.TimeStamp)
			fmt.Printf("\tHeight:%d\n", curBlock.Height)
			fmt.Printf("\tNonce:%d\n", curBlock.Nonce)
			fmt.Printf("\tTxs:%v\n", curBlock.Txs)
			for _, tx := range curBlock.Txs {
				fmt.Printf("\t\ttx-hash: %x\n", tx.TxHash)
				fmt.Printf("\t\t输入...\n")
				for _, vin := range tx.Vins {
					fmt.Printf("\t\t\ttvin-txHash: %x\n", vin.TxHash)
					fmt.Printf("\t\t\ttvin-vout: %x\n", vin.Vout)
					fmt.Printf("\t\t\ttvin-scriptSig: %x\n", vin.ScriptSig)
				}
				fmt.Printf("\t\t输出...\n")
				for _, vout := range tx.Vouts {
					fmt.Printf("\t\t\tvout=value:%d\n", vout.Value)
					fmt.Printf("\t\t\tvout-scriptPubkey:%s\n", vout.ScriptPubkey)
				}
			}
		} else {
			break
		}
	}
}

// MineNewBlock 实现挖矿功能：通过接收交易，生成区块
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	var block *Block
	// 搁置交易生成步骤
	var txs []*Transaction
	// 遍历交易参与者
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		// 生成新的交易
		tx := NewSimpleTransaction(address, to[index], value, bc, txs)
		// 追加到 txs 链表中
		txs = append(txs, tx)
	}

	// 从数据库中获取最新一个区块
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			// 获取最新区块的哈希值
			hash := b.Get([]byte("1"))
			// 获取最新区块
			blockBytes := b.Get(hash)
			// 反序列化
			block = Deserialize(blockBytes)
		}
		return nil
	})
	// 通过已拿到的区块生成新的区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	// 持久化新生成的区块到数据库中
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			err := b.Put(block.Hash, block.Serialize())
			if nil != err {
				log.Panicf("update the new block to db failed %v\n", err)
			}
			// 更新最新区块的哈希值
			err = b.Put([]byte("1"), block.Hash)
			if nil != err {
				log.Panicf("update the latest block hash to db failed %v\n", err)
			}
			bc.Tip = block.Hash
		}
		return nil
	})
}

// UnUTXOs 查找指定地址的 UTXO
/*
	遍历查找区块链数据库中的每一个区块中的每一个交易
	查找每一个交易中的每一个输出
	判断每个输出是否没满足下列条件
		1. 输入传入的地址
		2. 是否未被花费
			1. 遍历一次区块链数据库，将所有已花费的 OUTPUT 存入一个缓存
			2. 再次遍历区块链数据库，检查每一个 VOUT 是否包含在前面的已花费输出的缓存中
 */
func (bc *BlockChain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	// 1. 遍历数据库，查找与所有 address 相关的交易
	// 获取迭代器
	bcit := bc.Iterator()
	// 获取指定地址所有已花费输出
	spendTxOutputs := bc.SpentOutputs(address)
	// 当前地址的未花费输出列表
	var unUTXOS []*UTXO
	// 缓存迭代，查找缓存中的已花费输出
	for _, tx := range txs {
		// 判断 coinbaseTransaction
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {
				// 判断用户
				if in.CheckPubkeyWithAddress(address) {
					// 添加到已花费输出的 map 中
					key := hex.EncodeToString(in.TxHash)
					spendTxOutputs[key] = append(spendTxOutputs[key], in.Vout)
				}
			}
		}
	}
	// 优先遍历缓存中的 UTXO，如果余额足够，直接返回，如果不足，再遍历 db 文件中的 UTXO
	for _, tx := range txs {
		WorkCacheTx:
		for index, vout := range tx.Vouts {
			if vout.CheckPubkeyWithAddress(address) {
				if len(spendTxOutputs) != 0 {
					var isUtxoTx bool   // 判断交易是否被其他交易引用
					for txHash, indexArray := range spendTxOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if txHash == txHashStr {
							// 当前遍历到的交易已经有输出被其他交易的输入所引用
							isUtxoTx = true
							// 添加状态遍历，判断指定的 output 是否被引用
							var isSpentUTXO bool
							for _, voutIndex := range indexArray {
								if index == voutIndex {
									// 该输出被引用
									isSpentUTXO = true
									// 跳出当前 vout 判断逻辑，进行下一个输出判断
									continue WorkCacheTx
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, vout}
								unUTXOS = append(unUTXOS, utxo)
							}
						}
					}
					if isUtxoTx == false {
						// 说明当前交易中所有 address 相关的 outputs 都是 UTXO
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, vout}
					unUTXOS = append(unUTXOS, utxo)
				}
			}
		}
	}

	// 数据库迭代，不断获取下一个区块
	for {
		if !bcit.HasNext() {
			break
		}
		block := bcit.Next()
		// 遍历区块中的每笔交易
		for _, tx := range block.Txs {
			// 跳转
			work:
			for index, vout := range tx.Vouts {
				if vout.CheckPubkeyWithAddress(address) {
					// 当前 vout 输入传入地址
					if len(spendTxOutputs) != 0 {
						var isSpentOutput bool
						for txHash, indexArray := range spendTxOutputs {
							// txHash：当前输出所引用的交易哈希，indexArray：哈希关联的 vout 索引列表
							for _, i := range indexArray {
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									// txHash == hex.EncodeToString(tx.TxHash) 说明当前的交易 tx 至少已近有输出被其他交易的输入引用
									// index == i 说明当前输出正好被其他交易引用
									// 跳转到最外层循环，判断下一个 VOUT
									isSpentOutput = true
									continue work
								}
							}
						}
						if isSpentOutput == false {
							utxo := &UTXO{tx.TxHash, index, vout}
							unUTXOS = append(unUTXOS, utxo)
						}
					} else {
						// 将当前地址所有输出都添加到未花费输出中
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				}
			}
		}
	}
	return unUTXOS
}

// SpentOutputs 获取指定地址所有已花费输出
func (bc *BlockChain) SpentOutputs(address string) map[string][]int {
	// 已花费输出缓存
	spentTxOutputs := make(map[string][]int)
	// 获取迭代器对象
	bcit := bc.Iterator()
	for {
		if !bcit.HasNext() {
			break
		}
		block := bcit.Next()
		for _, tx := range block.Txs {
			// 排除 coinbase 交易
			if !tx.IsCoinbaseTransaction() {
				for _, in := range tx.Vins {
					if in.CheckPubkeyWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						// 添加到已花费输出的缓存中
						spentTxOutputs[key] = append(spentTxOutputs[key], in.Vout)
					}
				}
			}
		}
	}
	return spentTxOutputs
}

// 查询余额
func (bc *BlockChain) getBalance (address string) int {
	var amout int
	utxos := bc.UnUTXOs(address, []*Transaction{})
	for _, utxo := range utxos {
		amout += utxo.Output.Value
	}
	return amout
}

// FindSpendableUTXO 查找指定地址的可用 UTXO，超过 amount 就中断查找，更新当前数据库中指定地址的 UTXO 数量, txs：缓存中的交易列表（用于多比交易处理）
func (bc *BlockChain) FindSpendableUTXO(from string, amount int, txs []*Transaction) (int, map[string][]int) {
	// 可用的 UTXO
	spendableUTXO := make(map[string][]int)
	var value int
	utxos := bc.UnUTXOs(from, txs)
	// 遍历 UTXO
	for _, utxo := range utxos {
		value += utxo.Output.Value
		// 计算交易哈希
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}
	// 所有的循环遍历完成，仍然小于 amount，资金不足
	if value < amount {
		fmt.Printf("地址 [%s] 余额不足，当前余额 [%d]，转账金额 [%d]\n", from, value, amount)
		os.Exit(1)
	}
	return value, spendableUTXO
}