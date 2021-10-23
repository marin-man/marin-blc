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
	// 创建钱包
	fmt.Printf("\tcreatewallet -- 创建钱包\n")
	// 初始化区块链
	fmt.Printf("\tcreateblockchain -address address -- 创建区块链\n")
	// 打印完整的区块信息
	fmt.Printf("\tprintchain -- 输出区块信息\n")
	// 通过命令转账
	fmt.Printf("\tsend -from FROM -to TO -amount AMOUNT -- 发起转账\n")
	// 参数说明
	fmt.Printf("\t转账参数说明\n")
	fmt.Printf("\t\t-from FROM -- 转账源地址\n")
	fmt.Printf("\t\t-to TO -- 转账目标地址\n")
	fmt.Printf("\t\t-amount AMOUNT -- 转账金额\n")
	// 查询余额
	fmt.Printf("\tgetbalance -address FROM -- 查询指定地址的余额\n")
	fmt.Println("\t查询余额参数说明")
	fmt.Printf("\t\t-address -- 查询余额的地址")
	fmt.Printf("\taccounts -- 获取钱包地址列表\n")
	fmt.Printf("\tutxo -method METHOD -- 测试UTXO Table 功能中指定的方法\n")
	fmt.Printf("\t\tMETHOD -- 方法名\n")
	fmt.Printf("\t\t\treset -- 重置 UTXOtable\n")
	fmt.Printf("\t\t\tbalance -- 查找所有 UTXO\n")
	fmt.Printf("\tset_id -port PORT -- 设置节点号\n")
	fmt.Printf("\t\tPORT -- 访问的节点号")
	fmt.Printf("\tstart -- 启动节点服务")
}

// dbExit 判断数据库文件是否存在
func dbExit(nodeId string) bool {
	// 生成不同节点的数据库文件
	dbName := fmt.Sprintf(dbName, nodeId)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// 数据库文件不存在
		return false
	}
	return true
}

// BlockchainObject 获取一个 blockchain 对象
func BlockchainObject(nodeId string) *BlockChain {
	// 获取 DB
	dbName := fmt.Sprintf(dbName, nodeId)
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

// Run 命令行运行函数
func (cli *CLI) Run() {
	nodeId := GetEnvNodeId()
	// 检测参数数量
	IsValidArgs()
	// 新建相关命令
	// 添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// 输出区块链完整信息
	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 创建区块链
	createBLCWithGenesisBlockCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	// 发起交易
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	// 查询余额
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	// 创建钱包集合
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	// 获取地址列表
	getAccountsCmd := flag.NewFlagSet("accounts", flag.ExitOnError)
	// utxo 测试命令
	UTXOTestCmd := flag.NewFlagSet("utxo", flag.ExitOnError)
	// 节点号设置命令
	setNodeIdCmd := flag.NewFlagSet("set_id", flag.ExitOnError)
	// 节点服务启动命令
	startNodeCmd := flag.NewFlagSet("start", flag.ExitOnError)

	// 数据参数处理
	// 添加区块
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")
	// 创建区块时指定的矿工地址
	flagCreateBlockchainArg := createBLCWithGenesisBlockCmd.String("address", "troytan",
		"指定接收系统奖励的矿工地址")
	// 发起交易参数
	flagSendFromArg := sendCmd.String("from", "", "转账源地址")
	flagSendToArg := sendCmd.String("to", "", "转账目标地址")
	flagSendAmountArg := sendCmd.String("amount", "", "转账金额")
	// 查询余额命令行参数
	flagGetBalanceArg := getbalanceCmd.String("address", "", "要查询的地址")
	// UTXO 测试命令行参数
	flagUTXOArg := UTXOTestCmd.String("method", "", "UTXO Table 相关操作")
	// 端口号参数
	flagPortArg := setNodeIdCmd.String("port", "", "设置节点 ID")

	// 判断命令
	switch os.Args[1] {
	case "start" :
		if err := startNodeCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd start node server failed! %v\n", err)
		}
	case "set_id" :
		if err := setNodeIdCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd set node id failed! %v\n", err)
		}
	case "utxo" :
		if err :=UTXOTestCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd operate utxo table failed!%v\n", err)
		}
	case "accounts" :
		if err := getAccountsCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd get accounts failed! %v\n", err)
		}
	case "createwallet" :
		if err := createWalletCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd of create wallet failed! %v\n", err)
		}
	case "getbalance" :
		if err := getbalanceCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse cmd get balance failed %v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); nil != err {
			log.Panicf("parse send failed! %v\n", err)
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

	// 节点启动服务
	if startNodeCmd.Parsed() {
		cli.startNode(nodeId)
	}

	// 节点 ID 设置
	if setNodeIdCmd.Parsed() {
		if *flagPortArg == "" {
			fmt.Println("请输入端口号...")
			os.Exit(1)
		}
		cli.SetNodeId(*flagPortArg)
	}
	// utxo table 操作
	if UTXOTestCmd.Parsed() {
		switch *flagUTXOArg {
		case "balance":
			cli.TestFindUTXOMap()
		case "reset":
			cli.TestResetUTXO(nodeId)
		}
	}

	// 获取钱包地址列表
	if getAccountsCmd.Parsed() {
		cli.GetAccounts(nodeId)
	}

	// 创建钱包集合
	if createWalletCmd.Parsed() {
		cli.createWallets(nodeId)
	}

	// 查询余额
	if getbalanceCmd.Parsed() {
		if "" == *flagGetBalanceArg {
			fmt.Println("请输入查询地址...")
			os.Exit(1)
		}
		cli.getBalance(*flagGetBalanceArg, nodeId)
	}

	// 发起转账
	if sendCmd.Parsed() {
		if *flagSendFromArg == "" {
			fmt.Println("源地址不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Println(*flagSendFromArg)
		if *flagSendToArg == "" {
			fmt.Println("目标地址不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg == "" {
			fmt.Println("转账金额不能为空...")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n", JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n", JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]\n", JSONToSlice(*flagSendAmountArg))
		cli.send(JSONToSlice(*flagSendFromArg), JSONToSlice(*flagSendToArg), JSONToSlice(*flagSendAmountArg), nodeId)
	}

	// 输出区块链
	if printchainCmd.Parsed() {
		cli.printChain(nodeId)
	}
	// 创建区块链
	if createBLCWithGenesisBlockCmd.Parsed() {
		if *flagCreateBlockchainArg	== "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.createBlockchain(*flagAddBlockArg, nodeId)
	}
}