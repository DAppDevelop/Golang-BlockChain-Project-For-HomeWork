package BLC

import (
	"flag"
	"os"
	"log"
	"fmt"
)

type CLIYS struct{}

func (cli *CLIYS) RunYS() {

	isValidArgsYS()

	//配置./moac xxx 中xxx的命令参数
	//e.g. ./moac addblock
	addblockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	createblockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//关联命令参数
	flagAddBlockData := addblockCmd.String("data", "chenysh", "交易数据")
	flagCreateBlockchainWithCmd := createblockchainCmd.String("data", "GenesisBlock.......", "创世区块数据")

	switch os.Args[1] {
	case "addblock":
		//解析参数
		if err := addblockCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		if err := createblockchainCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	case "printchain":
		if err := printchainCmd.Parse(os.Args[2:]); err != nil {
			log.Panic(err)
		}
	default:
		printUsageYS()
		os.Exit(1)
	}

	//Parsed() -》是否执行过Parse()
	if addblockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsageYS()
			os.Exit(1)
		}

		cli.addBlockYS(*flagAddBlockData)
	}

	if createblockchainCmd.Parsed() {
		if *flagCreateBlockchainWithCmd == "" {
			printUsageYS()
			os.Exit(1)
		}

		cli.createGenesisBlockchainYS(*flagCreateBlockchainWithCmd)
	}

	if printchainCmd.Parsed() {
		cli.printchainYS()
	}

}

//输出使用指南
func printUsageYS() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -data -- 创世区块交易数据.")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")
}

func (cli *CLIYS) addBlockYS(data string) {
	if DBExistsYS() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObjectYS()

	defer blockchain.DBYS.Close()

	blockchain.AddBlockToBlockchainYS(data)

}

func (cli *CLIYS) printchainYS() {
	if DBExistsYS() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObjectYS()

	defer blockchain.DBYS.Close()

	blockchain.PrintchainYS()

}

func (cli *CLIYS) createGenesisBlockchainYS(data string) {
	CreatBlockchainWithGenesisBlockYS(data)
}

//判断参数是否有效
func isValidArgsYS() {
	if len(os.Args) < 2 {
		printUsageYS()
		os.Exit(1)
	}
}
