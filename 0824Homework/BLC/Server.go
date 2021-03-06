package BLC

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
)

/*
	启动服务器
 */
func startServerYS(nodeID string, mineAddress string) {
	//设置coinbase
	coinbaseAddressYS = mineAddress
	//拼接nodeID到ip后
	nodeAddressYS = fmt.Sprintf("localhost:%s", nodeID)
	//监听地址
	listener, err := net.Listen("tcp", nodeAddressYS)

	if err != nil {
		log.Panic(err)
	}

	defer listener.Close()

	bc := BlockchainObjectYS(nodeID)
	//defer bc.DB.Close()

	//判断是否为主节点, 非主节点的节点需要向主节点发送Version消息
	//fmt.Println(nodeAddress, knowNodes[0])
	if nodeAddressYS != knowNodesYS[0] {
		//fmt.Println("sendVersion")
		sendVersionYS(knowNodesYS[0], bc)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Panic(err)
		}

		fmt.Println("发送方已接入..", conn.RemoteAddr())

		go handleConnectionYS(conn, bc)
	}
}

/*
	处理请求结果
 */
func handleConnectionYS(conn net.Conn, bc *BlockchainYS) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}

	command := bytesToCommand(request[:COMMAND_LENGTHYS])

	fmt.Printf("接收到的命令是：%s\n", command)

	switch command {
	case COMMAND_VERSIONYS:
		handleVersionYS(request, bc)
	case COMMAND_GETBLOCKSYS:
		handleGetBlocksHashYS(request, bc)
	case COMMAND_INVYS:
		handleInvYS(request, bc)
	case COMMAND_GETDATAYS:
		handleGetDataYS(request, bc)
	case COMMAND_BLOCKDATAYS:
		handleGetBlockDataYS(request, bc)
	case COMMAND_TXSYS:
		handleTransactionsYS(request, bc)
	case COMMAND_REQUIREMINEYS:
		handleRequireMineYS(request, bc)
	case COMMAND_VERIFYBLOCKYS:
		handleVerifyBlockYS(request, bc)
	default:
		fmt.Println("无法识别....")
	}

	defer conn.Close()
}
