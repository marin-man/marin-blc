package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// 区块链迭代器管理文件

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

// Next 实现迭代函数 next
func (bcit *BlockChainIterator) Next() *Block {
	var block *Block
	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
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
	return block
}

// HasNext 查询是否有下一个值
func (bcit *BlockChainIterator) HasNext() bool {
	var block *Block
	var hasNext bool = true
	err := bcit.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			currentBlockBytes := b.Get(bcit.CurrentHash)
			block = Deserialize(currentBlockBytes)
			if nil == block.PrevBlockHash {
				hasNext = false
			}
		}
		return nil
	})
	if nil != err {
		log.Panicf("iterator the db failed %v\n", err)
	}
	return hasNext
}