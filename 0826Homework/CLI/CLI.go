package CLI

import (
	"flag"
	"os"
	"fmt"
	"encoding/json"
	"log"
	"chenyuanshan/0826Homework/RaftServer"
)

type CLIYS struct {}

func (cli *CLIYS)RunYS()  {

	setNodeCmd := flag.NewFlagSet("node", flag.ExitOnError)

	nameFrom := setNodeCmd.String("local", "", "本地节点地址")
	addressesFrom := setNodeCmd.String("addrs", "", "已知地址集")

	//设置一个参数 传入已知节点的Addr

	switch os.Args[1] {
	case "node":
		err := setNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if setNodeCmd.Parsed() {
		addrArray := JSONToArrayYS(*addressesFrom)
		if len(addrArray) == 0 || *nameFrom == "" {
			printUsage()
			os.Exit(1)
		}

		var addresses []RaftServer.AddrYS
		for _, ad := range addrArray {
			addr := ad
			address := RaftServer.AddrYS{addr}
			addresses = append(addresses, address)
		}
		raftServer := RaftServer.CreateNewRaftServerYS(*nameFrom, addresses)
		raftServer.RunYS()
	} else {
		fmt.Println("没解析参数")
		printUsage()
	}


}


// 标准的JSON字符串转数组
func JSONToArrayYS(jsonString string) []string {
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tnode -local 本地节点地址 -addrs 已知节点地址集 --启动节点")
}
