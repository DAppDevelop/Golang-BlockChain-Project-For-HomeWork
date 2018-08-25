package BLC

import (
	"fmt"
	"flag"
	"os"
	"log"
)

type CLIYS struct{}


func (cli *CLIYS) RunYS() {
	isValidArgsYS()

	//记录环境变量NODE_ID
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Println("没有设置NODE_ID")
		os.Exit(1)
	}

	fmt.Println("当前节点是:", nodeID)

	//1.---------创建flagset命令对象
	//e.g. ./moac addblock
	//./bc  命令 -参数名 参数
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	getAddresslistsCmd := flag.NewFlagSet("getaddresslists", flag.ExitOnError)
	createblockchainCmd := flag.NewFlagSet("create", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	resetCmd := flag.NewFlagSet("reset", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	setCoinbaseCmd := flag.NewFlagSet("coinbase", flag.ExitOnError)

	//2.----------设置命令后的参数对象
	flagFrom := sendCmd.String("from", "", "转账源地址")
	flagTo := sendCmd.String("to", "", "转账目的地址")
	flagAmount := sendCmd.String("amount", "", "转账金额")
	flagMine := sendCmd.String("mine", "", "本地挖矿")

	//createblockchainCmd 创世区块地址
	flagCoinbase := createblockchainCmd.String("address", "", "创世区块数据的地址")

	//getbalanceCmd
	flagGetbalanceWithAddress := getBalanceCmd.String("address", "", "要查询余额的账户.......")

	flagStartNodeWithMiner := startNodeCmd.String("miner", "", "挖矿奖励的地址")

	flagSetCoinbaseWithAddress := setCoinbaseCmd.String("address", "", "挖矿奖励地址")

	//3.----------解析参数
	switch os.Args[1] {
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "create":
		if err := createblockchainCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "print":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "createwallet":
		if err := createWalletCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "getaddresslists":
		if err := getAddresslistsCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "reset":
		if err := resetCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "startnode":
		if err := startNodeCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "coinbase":
		if err := setCoinbaseCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}

	default:
		printUsageYS()
		os.Exit(1)
	}

	//4.---------根据终端输入的命令执行对应的功能
	//Parsed() -》是否执行过Parse()
	if sendCmd.Parsed() {
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsageYS()
			os.Exit(1)
		}

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flagAmount)
		mine := true
		if *flagMine == "false" || *flagMine == "f" {
			mine = false
		}
		cli.sendYS(from, to, amount, nodeID, mine)
	}

	if createblockchainCmd.Parsed() {
		if *flagCoinbase == "" {
			fmt.Println("地址不能为空....")
			printUsageYS()
			os.Exit(1)
		}

		cli.createGenesisBlockchainYS(*flagCoinbase,nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printchainYS(nodeID)
	}

	if getBalanceCmd.Parsed() {
		if *flagGetbalanceWithAddress == "" {
			fmt.Println("地址不能为空....")
			printUsageYS()
			os.Exit(1)
		}

		cli.getBalanceYS(*flagGetbalanceWithAddress,nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.CreateWalletYS(nodeID)
	}

	if getAddresslistsCmd.Parsed() {
		cli.GetAddressListYS(nodeID)
	}

	if resetCmd.Parsed() {
		cli.ResetYS(nodeID)
	}

	if startNodeCmd.Parsed() {
		cli.startNodeYS(nodeID, *flagStartNodeWithMiner)
	}

	if setCoinbaseCmd.Parsed() {
		cli.setCoinbaseYS(nodeID, *flagSetCoinbaseWithAddress)
	}
}

/*
	输出使用指南
 */
func printUsageYS() {
	fmt.Println("Usage:")
	fmt.Println("\tcreatewallet -- 创建钱包")
	fmt.Println("\tgetaddresslists -- 获取所有的钱包地址")
	fmt.Println("\tcreate -address --创世区块交易数据.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -mine true/false  --转账交易")
	fmt.Println("\tprint --输出区块信息.")
	fmt.Println("\tgetbalance -address --获取address的余额.")
	fmt.Println("\treset --重置UTXOSet.")
	fmt.Println("\tstartnode --启动节点")
	fmt.Println("\tcoinbase -address --设置挖矿奖励地址")
}

/*
	判断参数是否有效
 */
func isValidArgsYS() {
	if len(os.Args) < 2 {
		printUsageYS()
		os.Exit(1)
	}
}