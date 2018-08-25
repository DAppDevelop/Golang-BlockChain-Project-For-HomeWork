package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"encoding/hex"
	"crypto/ecdsa"
	"bytes"
)

type BlockchainYS struct {
	TipYS []byte   //存储区块中最后一个块的hash值
	DBYS  *bolt.DB //对应的数据库对象
}

//1. 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlockYS(address string, nodeID string) {

	//设置dbname
	DBName := fmt.Sprintf(DBNameYS, nodeID) //"blockchain_3000.db"

	//判断数据库是否已经存
	if DBExistsYS(DBName) {
		fmt.Println("Genesis Block 已经存在...")
		os.Exit(1)
	}

	fmt.Println("创建创世区块....")

	//创建或打开数据库
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {

		//创建表
		b, err := tx.CreateBucketIfNotExists([]byte(BlockBucketNameYS))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建了一个coinbase Transaction
			txCoinbase := NewCoinbaseTransacionYS(address)
			// 创建创世区块
			genesisBlock := CreateGenesisBlockYS([]*TransactionYS{txCoinbase})

			//序列号block并存入数据库
			err := b.Put(genesisBlock.HashYS, gobEncode(genesisBlock))

			if err != nil {
				log.Panic(err)
			}

			//更新数据库最新区块hash
			err = b.Put([]byte("l"), genesisBlock.HashYS)

			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

/*
	处理交易命令数据, 输出对应的transactions
 */
func (blockchain *BlockchainYS)hanldeTransationsYS(from []string, to []string, amount []string, nodeId string) []*TransactionYS {
	var txs []*TransactionYS
	utxoSet := &UTXOSetYS{blockchain}

	for i := 0; i < len(from); i++ {
		//转换amount为int
		amountInt, _ := strconv.Atoi(amount[i])
		tx := NewSimpleTransation(from[i], to[i], int64(amountInt), utxoSet, txs, nodeId)
		//fmt.Println(tx)
		txs = append(txs, tx)
	}

	return txs
}


// 挖矿产生区块
func (blockchain *BlockchainYS) MineNewBlockYS(originalTxs []*TransactionYS) *BlockYS {
	/*
	奖励：reward：
	创建一个CoinBase交易--->Tx
	 */
	coinBaseTransaction := NewRewardTransacionYS()
	txs := []*TransactionYS{coinBaseTransaction}
	txs = append(txs, originalTxs...)

	//fmt.Println("交易的验证")
	//交易的验证：
	for _, tx := range txs {
		//fmt.Println(tx)
		//coinbase交易不验证
		if !tx.IsCoinBaseTransactionYS() {
			//fmt.Println(tx)
			if blockchain.VerifityTransactionYS(tx, txs) == false {
				log.Panic("数字签名验证失败。。。")
			}
		}
	}

	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	//获取最新的block
	var block BlockYS
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketNameYS))
		if b != nil {

			hash := b.Get([]byte("l"))

			blockBytes := b.Get(hash)

			gobDecode(blockBytes, &block)
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	//2. 根据数据库最新的block的信息,建立新的区块
	newBlock := NewBlockYS(txs, block.HeightYS+1, block.HashYS)
	//println(newBlock)

	return newBlock
}

func (blockchain *BlockchainYS)SaveNewBlockToBlockchainYS(newBlock *BlockYS)  {
	//将新区块存储到数据库

	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketNameYS))
		if b != nil {

			b.Put(newBlock.HashYS, gobEncode(newBlock))

			b.Put([]byte("l"), newBlock.HashYS)

			blockchain.TipYS = newBlock.HashYS

		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

//设计一个方法，用于获取指定用户的所有的未花费Txoutput
/*
UTXO模型：未花费的交易输出
	Unspent Transaction TxOutput
	txs --> 本次转账信息(查询时为nil)
 */
func (blc *BlockchainYS) UnSpentYS(address string, txs []*TransactionYS) []*UTXOYS {
	/*
	0.查询本次转账已经创建了的哪些transaction
	1.遍历数据库，获取每个block ---> Txs
	2.遍历所有交易：
		Inputs -- 记录为已花费
		Outputs -- 每个output
	 */

	//存储未花费的TxOutput
	var unSpentUTXOs [] *UTXOYS
	//存储已经花费的信息
	spentTxOutputMap := make(map[string][]int) // map[TxID] = []int{vout}

	//第一部分：先查询本次转账，已经产生了的Transanction
	for i := len(txs) - 1; i >= 0; i-- {
		unSpentUTXOs = caculateYS(txs[i], address, spentTxOutputMap, unSpentUTXOs)
	}

	for utxo := range unSpentUTXOs {
		fmt.Println("unSpentUTXO", utxo)
	}

	//第二部分：数据库里的Trasacntion
	it := blc.IteratorYS()
	for {

		//1、获取每个block
		block := it.NextYS()
		//2、遍历block的Txs
		//倒序遍历Transactions
		for i := len(block.TxsYS) - 1; i >= 0; i-- {
			unSpentUTXOs = caculateYS(block.TxsYS[i], address, spentTxOutputMap, unSpentUTXOs)
		}

		//3、判断退出
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PrevBlockHashYS)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}

	}

	return unSpentUTXOs
}

//计算对应address的未花费TXOutput
func caculateYS(tx *TransactionYS, address string, spentTxOutputMap map[string][]int, unSpentUTXOs []*UTXOYS) []*UTXOYS {
	//遍历每个tx：txID，Vins，Vouts

	//遍历所有的TxInput
	if !tx.IsCoinBaseTransactionYS() { //tx不是CoinBase交易，遍历TxInput
		for _, txInput := range tx.VinsYS {
			//txInput-->TxInput
			full_payload := Base58Decode([]byte(address))

			pubKeyHash := full_payload[1 : len(full_payload)-addressCheckSumLenYS]
			if txInput.UnlockWithAddressYS(pubKeyHash) {
				//txInput的解锁脚本(用户名) 如果和钥查询的余额的用户名相同，
				key := hex.EncodeToString(txInput.TxIDYS)
				spentTxOutputMap[key] = append(spentTxOutputMap[key], txInput.VoutYS)
				/*
				map[key]-->value TxInput.
				map[key] -->[]int
				 */
			}
		}
	}

	//遍历所有的TxOutput
outputs:
	for index, txOutput := range tx.VoutsYS { //index= 0,txoutput.锁定脚本：王二狗
		if txOutput.UnlockWithAddressYS(address) {
			if len(spentTxOutputMap) != 0 {
				var isSpentOutput bool //false
				//遍历map
				for txID, indexArray := range spentTxOutputMap { //143d,[]int{1}
					//遍历 记录已经花费的下标的数组
					for _, i := range indexArray {
						if i == index && hex.EncodeToString(tx.TxIDYS) == txID {
							isSpentOutput = true //标记当前的txOutput是已经花费
							continue outputs
						}
					}
				}

				if !isSpentOutput {
					//unSpentTxOutput = append(unSpentTxOutput, txOutput)
					//根据未花费的output，创建utxo对象--->数组
					utxo := &UTXOYS{tx.TxIDYS, index, txOutput}
					unSpentUTXOs = append(unSpentUTXOs, utxo)
				}

			} else {
				//如果map长度未0,证明还没有花费记录，output无需判断
				//unSpentTxOutput = append(unSpentTxOutput, txOutput)
				utxo := &UTXOYS{tx.TxIDYS, index, txOutput}
				unSpentUTXOs = append(unSpentUTXOs, utxo)
			}
		}
	}
	return unSpentUTXOs

}

/*
提供一个方法，返回用于一次转账的交易中，即将被使用为花费的utxo
 */
func (bc *BlockchainYS) FindSpentableUTXOsYS(from string, amount int64, txs []*TransactionYS) (int64, map[string][]int) {
	/*
	1.根据from获取到的所有的utxo
	2.遍历utxos，累加余额，判断，是否如果余额，大于等于要要转账的金额，


	返回：map[txID] -->[]int{下标1，下标2} --->Output
	 */
	var total int64

	spentableMap := make(map[string][]int)
	//1.获取所有的utxo ：10
	utxos := bc.UnSpentYS(from, txs)
	//2.找即将使用utxo：3个utxo
	for _, utxo := range utxos {
		total += utxo.OutputYS.ValueYS
		txIDstr := hex.EncodeToString(utxo.TxIDYS)
		spentableMap[txIDstr] = append(spentableMap[txIDstr], utxo.IndexYS)

		if total >= amount {
			break
		}
	}

	//3.判断total是否大于等于amount
	if total < amount {
		fmt.Printf("%s，余额不足，无法转账。。", from)
		os.Exit(1)
	}

	return total, spentableMap

}

//提供一个功能：查询余额
func (blc *BlockchainYS) GetBalanceYS(address string, txs []*TransactionYS) int64 {
	unSpentUTXOs := blc.UnSpentYS(address, txs)
	var total int64

	for _, utxo := range unSpentUTXOs {
		total += utxo.OutputYS.ValueYS
	}

	return total
}

// 遍历输出所有区块的信息
func (blc *BlockchainYS) PrintchainYS() {
	//创建迭代器
	blockIterator := blc.IteratorYS()

	for {
		block := blockIterator.NextYS()

		fmt.Println(block)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHashYS)

		//判断当期的block是否为创世区块（创世区块perblockhash为000000....）
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

/*
	迭代器
 */
func (blockchain *BlockchainYS) IteratorYS() *BlockchainIteratorYS {
	return &BlockchainIteratorYS{blockchain.TipYS, blockchain.DBYS}
}

/*
	判断数据库是否存在
 */
func DBExistsYS(DBName string) bool {
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		return false
	}
	return true
}

/*
	返回Blockchain对象
 */
func BlockchainObjectYS(nodeID string) *BlockchainYS {
	//因为已经知道数据库的名字，所以只要取出最新区块hash，既可以返回blockchain对象

	DBName := fmt.Sprintf(DBNameYS, nodeID)

	if DBExistsYS(DBName) {
		//fmt.Println("数据库已经存在。。。")
		//打开数据库
		db, err := bolt.Open(DBName, 0600, nil)
		if err != nil {
			log.Panic(err)
		}
		defer db.Close()

		var blockchain *BlockchainYS

		err = db.View(func(tx *bolt.Tx) error {
			//打开bucket，读取l对应的最新的hash
			b := tx.Bucket([]byte(BlockBucketNameYS))
			if b != nil {
				//读取最新hash
				hash := b.Get([]byte("l"))
				blockchain = &BlockchainYS{hash, db}
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		return blockchain
	} else {
		fmt.Println("数据库不存在，无法获取BlockChain对象。。。")
		return nil
	}
}

/*
	签名交易
 */
func (bc *BlockchainYS) SignTrasanctionYS(tx *TransactionYS, privateKey ecdsa.PrivateKey, txs [] *TransactionYS) {
	//签名：需要1,私钥，2.要签名的交易中的部分数据
	//1.判断要签名的tx，如果时coninbase交易直接返回
	if tx.IsCoinBaseTransactionYS() {
		return
	}

	//2.获取该tx中的Input，引用之前的transaction中的未花费的output，
	prevTxs := make(map[string]*TransactionYS)
	for _, input := range tx.VinsYS {
		txIDStr := hex.EncodeToString(input.TxIDYS)
		prevTxs[txIDStr] = bc.FindTransactionByTxIDYS(input.TxIDYS, txs)
	}

	//3.签名
	tx.SignYS(privateKey, prevTxs)
}

/*
	根据交易ID，获取对应的交易
 */
func (bc *BlockchainYS) FindTransactionByTxIDYS(txID []byte, txs [] *TransactionYS) *TransactionYS {
	//1.先查找未打包的txs
	for _, tx := range txs {
		if bytes.Compare(tx.TxIDYS, txID) == 0 {
			return tx
		}
	}
	//遍历数据库，获取blcok--->transaction
	iterator := bc.IteratorYS()
	for {
		block := iterator.NextYS()
		for _, tx := range block.TxsYS {
			if bytes.Compare(tx.TxIDYS, txID) == 0 {
				return tx
			}
		}

		//判断结束循环
		bigInt := new(big.Int)
		bigInt.SetBytes(block.PrevBlockHashYS)
		if big.NewInt(0).Cmp(bigInt) == 0 {
			break
		}
	}

	return &TransactionYS{}
}

/*
	验证交易的数字签名
 */
func (bc *BlockchainYS) VerifityTransactionYS(tx *TransactionYS, txs []*TransactionYS) bool {
	//要想验证数字签名：私钥+数据 (tx的副本+之前的交易)
	//2.获取该tx中的Input，引用之前的transaction中的未花费的output
	prevTxs := make(map[string]*TransactionYS)
	for _, input := range tx.VinsYS {
		txIDStr := hex.EncodeToString(input.TxIDYS)
		prevTxs[txIDStr] = bc.FindTransactionByTxIDYS(input.TxIDYS, txs)
	}

	if len(prevTxs) == 0 {
		fmt.Println("没找到对应交易")
	} else {
		//fmt.Println("preTxs___________________________________")
		//fmt.Println(prevTxs)
	}

	//验证
	return tx.VerifityYS(prevTxs)
	//return true
}

/*
	获取所有区块中的UTXO
	map[string]*TxOutputs  交易id-->[]*UTXO (这笔交易下的UTXO集合)
*/
func (bc *BlockchainYS) FindUnspentUTXOMapYS() map[string]*TxOutputsYS {

	iterator := bc.IteratorYS()

	utxoMap := make(map[string]*TxOutputsYS)

	//已花费的input map
	spentedMp := make(map[string][]*TXInputYS)

	//遍历所有block
	for {
		block := iterator.NextYS()

		//倒序遍历block里面的TXs
		for i := len(block.TxsYS) - 1; i >= 0; i-- {
			//收集input
			tx := block.TxsYS[i]                     //当期的TX交易
			txIDStr := hex.EncodeToString(tx.TxIDYS) //TXID string

			txOutputs := &TxOutputsYS{[]*UTXOYS{}}

			//coinbase不处理Vins
			if !tx.IsCoinBaseTransactionYS() {
				for _, txInput := range tx.VinsYS {
					txIDStr := hex.EncodeToString(txInput.TxIDYS)
					spentedMp[txIDStr] = append(spentedMp[txIDStr], txInput)
				}
			}

			//根据spentedMp,遍历outputs 找出 UTXO
		outputLoop:
			for index, txOutput := range tx.VoutsYS {

				if len(spentedMp) > 0 {
					//isSpent := false
					inputs := spentedMp[txIDStr] //如果inputs 存在, 则对应的交易里面某笔output肯定已经被消费
					for _, input := range inputs {
						//判断input对应的是否当期的output
						if index == input.VoutYS && input.UnlockWithAddressYS(txOutput.PubKeyHashYS) {
							//此笔output已被消费
							//isSpent = true
							continue outputLoop
						}
					}

					//if isSpent == false {
					//outputs 加进utxoMap
					utxo := &UTXOYS{tx.TxIDYS, index, txOutput}
					txOutputs.UTXOsYS = append(txOutputs.UTXOsYS, utxo)
					//}
				} else {
					//outputs 加进utxoMap
					utxo := &UTXOYS{tx.TxIDYS, index, txOutput}
					txOutputs.UTXOsYS = append(txOutputs.UTXOsYS, utxo)
				}
			}

			if len(txOutputs.UTXOsYS) > 0 {
				utxoMap[txIDStr] = txOutputs
			}

		}

		//退出条件
		hashBigInt := new(big.Int)
		hashBigInt.SetBytes(block.PrevBlockHashYS)
		if big.NewInt(0).Cmp(hashBigInt) == 0 {
			break
		}
	}

	return utxoMap
}

/*
	获取blockchain最高高度
 */
func (bc *BlockchainYS) GetBestHeightYS() int64 {
	bestBlockChain := bc.IteratorYS().NextYS()
	return bestBlockChain.HeightYS
}

/*
	返回blockchain里面所有block的hash
 */
func (bc *BlockchainYS) getBlocksHashesYS() [][]byte {
	//迭代
	iterator := bc.IteratorYS()

	var blocksHashes [][]byte

	for {
		block := iterator.NextYS()

		blocksHashes = append(blocksHashes, block.HashYS)

		bigInt := new(big.Int)
		bigInt.SetBytes(block.PrevBlockHashYS)

		if big.NewInt(0).Cmp(bigInt) == 0 {
			break
		}
	}

	return blocksHashes
}

/*
	根据hash,获取对应的block
	hash --> Block
 */
func (bc *BlockchainYS) GetBlockByHashYS(hash []byte) *BlockYS {
	var block BlockYS

	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()


	//遍历
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketNameYS))
		if b != nil {
			blockBytes := b.Get(hash)
			//block = DeserializeBlock(blockBytes)
			gobDecode(blockBytes, &block)
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &block
}

/*
	添加一个block到blockChain里面
 */
func (bc *BlockchainYS) AddBlockYS(block *BlockYS) {
	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketNameYS))
		if b != nil {
			//判断次区块是否已经在本地的数据库里面
			blockBytes := b.Get(block.HashYS)
			if blockBytes != nil {
				return nil
			}

			err := b.Put(block.HashYS, gobEncode(block))
			if err != nil {
				log.Panic(err)
			}

			//判断新添加的block高度是否比当期最高高度高,是的话替换l
			lastBlockHash := b.Get([]byte("l"))
			lastBlockBytes := b.Get(lastBlockHash)
			//lastBlock := DeserializeBlock(lastBlockBytes)
			var lastBlock BlockYS
			gobDecode(lastBlockBytes, &lastBlock)

			if lastBlock.HeightYS < block.HeightYS {
				b.Put([]byte("l"), block.HashYS)
				bc.TipYS = block.HashYS
			}

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
