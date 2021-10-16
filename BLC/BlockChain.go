package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
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
func CreateBlockChain() *BlockChain{
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
			// 创建一个创世块
			genesisBlock := CreateGenesisBlock([]byte("init"))
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
func (bc *BlockChain) AddBlock(data []byte) {
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
			newBlock := NewBlock(latestBlock.Height + 1, latestBlock.Hash, data)
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
			curBlock := bcit.Next()
			fmt.Printf("Hash: %x\nPrevBlockHash: %x\nTimeStamp: %x\nData: %v\nHeight: %d\nNonce: %d\n",
				curBlock.Hash, curBlock.PrevBlockHash, curBlock.TimeStamp, curBlock.Data, curBlock.Height, curBlock.Nonce)
		} else {
			break
		}
	}
}