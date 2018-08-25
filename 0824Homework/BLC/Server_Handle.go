package BLC

import (
	"encoding/gob"
	"bytes"
	"log"
	"fmt"
	"os"
)

/*
	处理version命令
	1.根据本地区块高度以及版本信息判断后续操作
	本地高度>对方高度 -> 向对方发送本地的version命令消息
	对方高度>本地高度 -> 向对方请求对方的区块链信息
 */
func handleVersionYS(request []byte, bc *BlockchainYS) {

	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	var version VersionYS

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&version)
	if err != nil {
		log.Panic(err)
	}

	//3.操作bc，获取自己的最后block的height
	height := bc.GetBestHeightYS()
	foreignerBestHeight := version.BestHeightYS

	//4.根对方的比较, 相同则不做操作
	if height > foreignerBestHeight {
		//当前节点比对方节点高度高
		sendVersionYS(version.AddrFromYS, bc)
	} else if foreignerBestHeight > height {
		//当前节点比对方节点高度低,向对方节点请求对方节点的blockchain hash集
		sendGetBlocksHashYS(version.AddrFromYS)
	}

}

/*
	处理getblocks命令
	向对方发送本地的区块链hash集
 */
func handleGetBlocksHashYS(request []byte, bc *BlockchainYS) {
	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	var getblocks GetBlocks

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&getblocks)
	if err != nil {
		log.Panic(err)
	}

	blocksHashes := bc.getBlocksHashesYS()

	sendInvYS(getblocks.AddrFrom, BLOCK_TYPE, blocksHashes)
}

/*
	处理Inv命令
	1. block type :  如果本地区块

 */
func handleInvYS(request []byte, bc *BlockchainYS) {
	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	var inv InvYS

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&inv)
	if err != nil {
		log.Panic(err)
	}

	if inv.TypeYS == BLOCK_TYPE {
		//获取hashes中第一个hash,请求对方返回此hash对应的block
		hash := inv.ItemsYS[0]
		sendGetDataYS(inv.AddrFromYS, BLOCK_TYPE, hash)

		//保存items剩余未请求的hashes到变量blockArray(handleBlockData 方法会用到)
		if len(inv.ItemsYS) > 0 {
			blockArrayYS = inv.ItemsYS[1:]
		}

	} else if inv.TypeYS == TX_TYPE {

	}
}

func handleGetDataYS(request []byte, bc *BlockchainYS) {
	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	var getData GetData

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&getData)
	if err != nil {
		log.Panic(err)
	}

	if getData.Type == BLOCK_TYPE {
		block := bc.GetBlockByHashYS(getData.Hash)
		sendBlockYS(getData.AddrFrom, block)
	} else if getData.Type == TX_TYPE {

	}
}

func handleGetBlockDataYS(request []byte, bc *BlockchainYS) {
	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	var getBlockData BlockData

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&getBlockData)
	if err != nil {
		log.Panic(err)
	}

	blockBytes := getBlockData.Block
	//block := DeserializeBlock(blockBytes)
	var block BlockYS
	gobDecode(blockBytes, &block)
	//fmt.Println(&block)
	bc.AddBlockYS(&block)

	if len(blockArrayYS) == 0 {
		utxoSet := UTXOSetYS{bc}
		utxoSet.ResetUTXOSetYS()

	}

	if len(blockArrayYS) > 0 {
		hash := blockArrayYS[0]
		sendGetDataYS(getBlockData.AddrFrom, BLOCK_TYPE, hash)
		blockArrayYS = blockArrayYS[1:]
	}

}

/*
	主节点处理接收到的交易
 */
func handleTransactionsYS(request []byte, bc *BlockchainYS) {
	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	var txs []*TransactionYS

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&txs)
	if err != nil {
		log.Panic(err)
	}

	//发送到挖矿节点
	sendTransactionToMinerYS(knowNodesYS[1], txs)

	//for _, tx := range txs {
	//	//fmt.Println("处理获取到的txs")
	//	//fmt.Println(tx)
	//}
}

func handleRequireMineYS(request []byte, bc *BlockchainYS) {
	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]
	//fmt.Println("反序列化得到的txbytes：")
	//fmt.Printf("%x",commandBytes)
	//fmt.Println("-----")

	//2.反序列化--->version
	var txs []*TransactionYS

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&txs)
	if err != nil {
		log.Panic(err)
	}

	//fmt.Printf("%x",gobEncode(txs))

	nodeID := os.Getenv("NODE_ID")
	txp := NewTXPoolYS(nodeID)
	//将txs保存到交易池
	txp.TxsYS = append(txp.TxsYS, txs...)
	//for _, tx := range txp.Txs {
	//	fmt.Println(tx)
	//}
	txp.saveFileYS(nodeID)

	const packageNum = 1

	//2. 判断交易池是否有足够的交易
	if len(txp.TxsYS) > 0 {
		//开始挖矿
		fmt.Println("开始挖矿")

		blockchain := BlockchainObject(nodeID)

		//取出要打包的交易
		//packageTx := txp.Txs[:packageNum]
		newBlock := blockchain.MineNewBlockYS(txs)
		//fmt.Println(newBlock)
		txp.TxsYS = txp.TxsYS[packageNum:]
		txp.saveFileYS(nodeID)
		//发送newBlock 给主节点验证工作量证明
		sendNewBlockToMainYS(knowNodesYS[0], newBlock)
	}
}

func handleVerifyBlockYS(request []byte, blockchain *BlockchainYS) {
	//1.从request中获取版本的数据：[]byte
	commandBytes := request[COMMAND_LENGTH:]

	//2.反序列化--->version
	var block *BlockYS

	decoder := gob.NewDecoder(bytes.NewReader(commandBytes))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	pow := NewProofOfWork(block)
	if pow.IsValid() {
		blockchain.SaveNewBlockToBlockchainYS(block)
		utxoSet := &UTXOSetYS{blockchain}
		utxoSet.UpdateYS()

		//这里直接调起一次version命令  更新挖矿节点的区块
		sendVersionYS(knowNodesYS[1], blockchain)
	}

}
