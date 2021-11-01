package network

import (
	"bkc/core"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

// 请求处理文件管理

// handleVersion version
func handleVersion(request []byte, bc *core.BlockChain) {
	fmt.Println("the request of version handle...")
	var buffer bytes.Buffer
	var data Version
	// 1. 解析请求
	dataBytes := request[12:]
	// 2. 生成 version 结构
	buffer.Write(dataBytes)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the version struct failed! %v\n", err)
	}
	// 3. 获取请求方的区块高度
	versionHeight := data.Height
	// 4. 获取自身节点的区块高度
	height := bc.GetHeight()
	fmt.Printf("height : %v, versionHeigth : %v\n", height, versionHeight)
	if height > int64(versionHeight) {
		// 如果当前节点的区块高度大于 versionHeight，将当前节点版本信息发送给请求节点
		sendVersion(data.AddrFrom, bc)
	} else if height < int64(versionHeight) {
		// 如果当前接待你区块高度小于 versionHeight，向发送方发起同步数据的请求
		sendGetBlocks(data.AddrFrom)
	}
}

// handleGetBlocks 数据同步请求处理
func handleGetBlocks(request []byte, bc *core.BlockChain) {
	fmt.Println("the request of get blocks handle...")
	var buffer bytes.Buffer
	var data GetBlocks
	// 1. 解析请求
	dataBytes := request[12:]
	// 2. 生成 getblocks 结构
	buffer.Write(dataBytes)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the getblocks struct failed! %v\n", err)
	}
	// 3. 获取区块链所有区块的哈希
	hashes := bc.GetBlockHashes()
	sendInv(data.AddrFrom, hashes)

}

// handleInv
func handleInv(request []byte, bc *core.BlockChain) {
	fmt.Println("the request of inv handle...")
	var buffer bytes.Buffer
	var data Inv
	// 1. 解析请求
	dataBytes := request[12:]
	// 2. 生成 inv 结构
	buffer.Write(dataBytes)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the inv struct failed! %v\n", err)
	}
	for _, hash := range data.Hashes {
		sendGetData(data.AddrFrom, hash)
	}
}

// handleGetData 处理获取指定区块的请求
func handleGetData(request []byte, bc *core.BlockChain) {
	fmt.Println("the request of get block handle...")
	var buffer bytes.Buffer
	var data GetData
	// 1. 解析请求
	dataBytes := request[12:]
	// 2. 生成 getData 结构
	buffer.Write(dataBytes)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the getData struct failed! %v\n", err)
	}
	// 3. 通过传过来的区块哈希，获取本地节点的区块
	blockBytes := bc.GetBlock(data.ID)
	sendBlock(data.AddrFrom, blockBytes)
}

// handleBlock 接收到新区块时，进行处理
func handleBlock(request []byte, bc *core.BlockChain) {
	fmt.Println("the request of handle block handle...")
	var buffer bytes.Buffer
	var data BlockData
	// 1. 解析请求
	dataBytes := request[12:]
	// 2. 生成 getData 结构
	buffer.Write(dataBytes)
	decoder := gob.NewDecoder(&buffer)
	if err := decoder.Decode(&data); nil != err {
		log.Panicf("decode the blockData struct failed! %v\n", err)
	}
	// 3. 将接收到的区块添加到区块链中
	blockBytes := data.Block
	block := core.Deserialize(blockBytes)
	bc.AddBlock(block)
	// 4. 更新 UTXO
	utxoSet := core.UTXOSet{bc}
	utxoSet.Update()
}