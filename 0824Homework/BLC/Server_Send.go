package BLC

import (
	"net"
	"log"
	"io"
	"bytes"
)

/*
	所有消息都是通过这个方法来发送到其他节点
 */
func sendDataYS(to string, data []byte) {
	//fmt.Println("向",to,"发送",data)
	conn, err := net.Dial("tcp", to)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	//发送数据
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

/*
	发送本地版本/区块高度
 */
func sendVersionYS(to string, bc *BlockchainYS) {
	//1.创建对象
	bestHeight := bc.GetBestHeightYS()//获取当前节点区块链高度
	version := &VersionYS{NODE_VERSIONYS, bestHeight, nodeAddressYS}

	sendCommandDataYS(COMMAND_VERSIONYS, version, to)
}

/*
	发送请求要获取对方blockhash的消息
 */
func sendGetBlocksHashYS(to string) {
	//1.创建对象
	getBlocks := GetBlocksYS{nodeAddressYS}

	sendCommandDataYS(COMMAND_GETBLOCKSYS, getBlocks, to)
}

/*
	发送所有blockHash 数组的消息
 */
func sendInvYS(to string, kind string, data [][]byte) {
	//1.创建对象
	inv := InvYS{nodeAddressYS, kind, data}

	sendCommandDataYS(COMMAND_INVYS, inv, to)
}

/*
	发送请求对方根据hash返回对应的block的消息
 */
func sendGetDataYS(to string, kind string, hash []byte) {
	//1.创建对象
	getData := GetDataYS{nodeAddressYS, kind, hash}

	sendCommandDataYS(COMMAND_GETDATAYS, getData, to)
}

/*
	发送block对象给对方
 */
func sendBlockYS(to string, block *BlockYS) {
	//1.创建对象
	blockData := BlockDataYS{nodeAddressYS, gobEncode(block)}

	sendCommandDataYS(COMMAND_BLOCKDATAYS, blockData, to)
}

/*
	发送交易信息到主节点
 */
func sendTransactionToMainNodeYS(to string, txs []*TransactionYS)  {
	sendCommandDataYS(COMMAND_TXSYS, txs, to)
}

func sendTransactionToMinerYS(to string, txs []*TransactionYS)  {
	sendCommandDataYS(COMMAND_REQUIREMINEYS, txs, to)
}

func sendNewBlockToMainYS(to string, block *BlockYS) {
	sendCommandDataYS(COMMAND_VERIFYBLOCKYS, block, to)
}



func sendCommandDataYS(command string, data interface{}, to string)  {
	//2.对象序列化为[]byte
	payload := gobEncode(data)
	//3.拼接命令和对象序列化
	request := append(commandToBytes(command), payload...)
	//4.发送消息
	sendDataYS(to, request)
}


