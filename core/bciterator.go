package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

// 区块链迭代器管理文件
var hasNext = true

// BlockChainIterator 迭代器基本结构
type BlockChainIterator struct {
	DB				*bolt.DB	// 迭代目标
	CurrentHash		[]byte		// 当前迭代目标的哈希
}

// Iterator 创建迭代器对象
func (blc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		DB: blc.DB,
		CurrentHash: blc.Tip,
	}
}

// PreBlock 返回当前区块数据并更新当前区块哈希
func (bcit *BlockChainIterator) PreBlock() (*Block, bool) {
	var block *Block
	// 根据 hash 获取块数据
	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if nil != b {
			currentBlockBytes := b.Get(bcit.CurrentHash)
			block = Deserialize(currentBlockBytes)
			// 更新迭代器中的哈希值
			bcit.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if nil != err {
		log.Panicf("iterator the db failed %v\n", err)
	}
	// 返回区块
	return block, len(bcit.CurrentHash) > 0
}

// DBExits 判断区块链是否已经存在
func DBExits(nodeId string) bool {
	dbName := fmt.Sprintf(DBName, nodeId)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}