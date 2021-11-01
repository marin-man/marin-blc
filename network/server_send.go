package network

import (
	"bkc/core"
	"bkc/utils"
	"bytes"
	"io"
	"log"
	"net"
)

// sendMessage 发送请求
func sendMessage(to string, message []byte) {
	// 1. 连接上服务器
	conn, err := net.Dial(PROTOCOL, to)
	if nil != err {
		log.Panicf("connect to server [%s] failed! %v\n", to, err)
	}
	defer conn.Close()
	// 要发送的数据
	_, err = io.Copy(conn, bytes.NewReader(message))
	if nil != err {
		log.Panicf("add the data to conn failed! %v", err)
	}
}

// sendVersion 区块链版本验证
func sendVersion(toAddress string, bc *core.BlockChain) {
	// 1. 获取当前节点的区块高度
	height := bc.GetHeight()
	// 2. 组装生成 version
	versionData := Version{Height: int(height), AddrFrom: nodeAddress}
	// 3. 组装成要发送的请求
	data := utils.GobEncode(versionData)
	// 4. 将命令与版本组装成完整的请求
	request := append(CommandToBytes(CMD_VERSION), data...)
	// 5. 发送请求
	sendMessage(toAddress, request)
}

// sendGetBlocks 从指定节点同步数据
func sendGetBlocks(toAddress string) {
	// 1. 生成数据
	data := utils.GobEncode(GetBlocks{AddrFrom: nodeAddress})
	// 2. 组装请求
	request := append(CommandToBytes(CMD_GETBLOCKS), data...)
	// 3. 发送请求
	sendMessage(toAddress, request)
}

// sendGetData 发送获取指定节点请求
func sendGetData(toAddress string, hash []byte) {
	// 1. 生成数据
	data := utils.GobEncode(GetData{AddrFrom: nodeAddress, ID: hash})
	// 2. 组装请求
	request := append(CommandToBytes(CMD_GETDATA), data...)
	// 3. 发送请求
	sendMessage(toAddress, request)
}

// sendInv 向其他节点展示
func sendInv(toAddress string, hashes [][]byte) {
	// 1. 生成数据
	data := utils.GobEncode(Inv{AddrFrom: nodeAddress, Hashes: hashes})
	// 2. 组装请求
	request := append(CommandToBytes(CMD_INV), data...)
	// 3. 发送请求
	sendMessage(toAddress, request)
}

// sendBlock 发送区块信息
func sendBlock(toAddress string, block []byte)  {
	// 1. 生成数据
	data := utils.GobEncode(BlockData{AddrFrom: nodeAddress, Block: block})
	// 2. 组装请求
	request := append(CommandToBytes(CMD_BLOCK), data...)
	// 3. 发送请求
	sendMessage(toAddress, request)
}