package core

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 钱包集合管理文件

// 钱包集合持久化文件
const walletFile = "Wallets_%s.dat"

// Wallets 钱包集合的基本结构
type Wallets struct {
	Wallets map[string] *Wallet // key:地址  value:钱包结构
}

// NewWallets 初始化钱包集合
func NewWallets(nodeId string) *Wallets {
	// 从钱包文件中获取钱包信息
	walletFile := fmt.Sprintf(walletFile, nodeId)
	// 1. 判断文件是否存在
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.Wallets = make(map[string] *Wallet)
		return wallets
	}
	// 2. 文件存在，读取内容
	fileContent, err := ioutil.ReadFile(walletFile)
	if nil != err {
		log.Panicf("read the file content failed %v\n", err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if nil != err {
		log.Panicf("decode the file content failed! %v\n", err)
	}
	return &wallets
}

// CreateWallet 添加新的钱包到集合中
func (wallets *Wallets) CreateWallet(nodeId string)  {
	// 1. 创建钱包
	wallet := NewWallet()
	// 2. 添加
	wallets.Wallets[string(wallet.GetAddress())] = wallet
	// 3. 持久化钱包信息
	wallets.SaveWallets(nodeId)
}

// SaveWallets 持久化钱包信息(存储到文件中)
func (wallets *Wallets) SaveWallets(nodeId string) {
	walletFile := fmt.Sprintf(walletFile, nodeId)
	var content bytes.Buffer	// 钱包内容
	gob.Register(elliptic.P256())   // 注册256椭圆，注册之后，可以直接在内部对 curve 的接口进行编码
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&wallets)
	if nil != err {
		log.Panicf("encode the struct of wallets failed %v\n", err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if nil != err {
		log.Panicf("write the content of wallet into file [%s] failed! %v\n", walletFile, err)
	}
}