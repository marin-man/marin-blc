package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 共识算法管理文件

// 实现 POW 实例以及相关功能

// 目标值难度
const targetBit = 16

// ProofOfWork 工作量证明的结构
type ProofOfWork struct {
	Block 	*Block		// 需要共识验证的区块
	target	*big.Int	// 目标难度的哈希（大数据存储）
}

// NewProofOfWork 创建一个 POW 对象
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// 假设数据总长度为 8 位，需求：需要满足前两位为 0，才能解决问题
	target = target.Lsh(target, 256 - targetBit)
	// strTarget
	return &ProofOfWork{
		Block: block,
		target: target,
	}
}

// Run 执行 pow，比较哈希值，返回哈希值以及碰撞的次数
func (pow *ProofOfWork) Run() ([]byte, int) {
	// 碰撞次数
	var hashInt big.Int
	var hash [32]byte    // 生成的哈希值
	nonce := 0
	// 无限循环，生成符合调整的哈希值
	for {
		// 生成准备数据
		dataBytes := pow.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		// 检测生成的哈希值是否符合条件
		if pow.target.Cmp(&hashInt) == 1 {
			// 找到了符合条件的哈希值，中断循环
			break
		}
		nonce++
	}
	fmt.Printf("\n碰撞次数：%v\n", nonce)
	return hash[:], nonce
}

// prepareData 生成准备数据
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	// 拼接而区块属性，进行哈希计算
	timeStampBytes := IntToHex(pow.Block.TimeStamp)
	heightBytes := IntToHex(pow.Block.Height)
	data := bytes.Join([][]byte{
		heightBytes,
		timeStampBytes,
		pow.Block.PrevBlockHash,
		pow.Block.HashTransaction(),
		IntToHex(targetBit),
		IntToHex(nonce),
	}, []byte{})
	return data
}