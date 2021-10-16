package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// 区块基本结构与功能管理文件

// Block 实现一个最基本的区块结构
type Block struct {
	TimeStamp     	int64  	// 区块时间
	Hash          	[]byte 	// 当前区块哈希
	PrevBlockHash 	[]byte 	// 前区块哈希
	Height       	int64  	// 区块高度
	Data          	[]byte 	// 交易数据
	Nonce			int64	// 在运行 pow 时生成的哈希值，也代表 pow 运行时动态修改的数据
}

// NewBlock 新建区块
func NewBlock(Heigth int64, prevBlockHash []byte, data []byte) *Block {
	block := Block{
		TimeStamp:     time.Now().Unix(),
		Hash:          nil,
		PrevBlockHash: prevBlockHash,
		Height:        Heigth,
		Data:          data,
	}
	pow := NewProofOfWork(&block)
	// 执行工作量证明算法
	hash, nonce := pow.Run()
	// 生成当前区块哈希
	block.Hash = hash
	block.Nonce = int64(nonce)
	return &block
}

// CreateGenesisBlock 生成创世块
func CreateGenesisBlock(data []byte) *Block {
	return NewBlock(1, nil, data)
}

// Serialize 区块结构序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	// 新建编码对象
	encoder := gob.NewEncoder(&buffer)
	// 编码（序列化）
	err := encoder.Encode(block)
	if  nil != err {
		log.Panicln("serialize the block to []byte failed %v\n", err)
	}
	return buffer.Bytes()
}

// Deserialize 区块结构反序列化
func Deserialize(blockBytes []byte) *Block {
	var block Block
	// 新建 decoder 对象
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if nil != err {
		log.Panicf("deserialize the []byte to block failed %v\n", err)
	}
	return &block
}