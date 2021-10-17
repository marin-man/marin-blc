package BLC

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

// 对 blockchain 的命令行操作进行管理

// CLI client 对象
type CLI struct {
}

// PrintUsage 用法展示
func PrintUsage()  {
	fmt.Println("Usage:")
	// 初始化区块链
	fmt.Printf("createblockchain -- 创建区块链\n")
	// 添加区块
	fmt.Printf("addblock -data DATA-- 添加区块\n")
	// 打印完整的区块信息
	fmt.Printf("printchain -- 输出区块信息\n")
}

// createBlockchain 初始化区块链
func (cli *CLI) createBlockchain() {
	CreateBlockChain()
}

// addBlock 添加区块
func (cli *CLI) addBlock(data string) {
	if !dbExit() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	blockchain.AddBlock([]byte(data))
}

// printChain 打印完整区块链信息
func (cli *CLI) printChain() {
	if !dbExit() {
		fmt.Println("数据库不存在...")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	blockchain.PrintChain()
}

// dbExit 判断数据库文件是否存在
func dbExit() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// 数据库文件不存在
		return false
	}
	return true
}

// BlockchainObject 获取一个 blockchain 对象
func BlockchainObject() *BlockChain {
	// 获取 DB
	db, err := bolt.Open(dbName, 0600, nil)
	if nil != err {
		log.Panicf("open the db [%s] failed! %v\n", dbName, err)
	}
	// 获取 Tip
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if nil != b {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	if nil != err {
		log.Panicf("get the blockchain object failed ! %v\n", err)
	}
	return &BlockChain{
		DB: db,
		Tip: tip,
	}
}

// IsValidArgs 参数数量检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		// 直接退出
		os.Exit(1)
	}
}

// Run 命令行运行函数
func (cli *CLI) Run() {
	// 检测参数数量
	IsValidArgs()
	// 新建相关命令
	// 添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// 输出区块链完整信息
	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 创建区块链
	createBLCWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	// 数据参数处理
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")
	// 判断命令
	switch os.Args[1] {
	case "addblock" :
		if err := addBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse addBlockCmd failed! %v\n", err)
		}
	case "printchain" :
		if err := printchainCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createBLCWithGenesisBlockCmd failed! %v\n", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlockCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse createBLCWithGenesisBlockCmd failed %v\n", err)
		}
	default:
		// 没有传递任何命令或者传递命令不在上面的命令列表当中
		PrintUsage()
		os.Exit(1)
	}

	// 添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.addBlock(*flagAddBlockArg)
	}

	// 输出区块链
	if printchainCmd.Parsed() {
		cli.printChain()
	}
	// 创建区块链
	if createBLCWithGenesisBlockCmd.Parsed() {
		cli.createBlockchain()
	}
}