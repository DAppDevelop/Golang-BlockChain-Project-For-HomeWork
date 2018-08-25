package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"os"
	"bytes"
)

type UTXOSetYS struct {
	blockChainYS *BlockchainYS
}

/*
	查询block块中所有的未花费utxo：执行FindUnspentUTXOMap--->map
 */
func (utxoset *UTXOSetYS) ResetUTXOSetYS() {
	//blockchain.FindUnspentUTXOMap 涉及到数据库, 现在要放在下面打开数据库的命令之前.不然卡死
	utxoMap := utxoset.blockChainYS.FindUnspentUTXOMapYS()

	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	fmt.Println(DBName)
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()


	err = db.Update(func(tx *bolt.Tx) error {
		//1.utxoset表存在，删除
		b := tx.Bucket([]byte(UTXOSetBucketNameYS))
		if b != nil {
			err := tx.DeleteBucket([]byte(UTXOSetBucketNameYS))
			if err != nil {
				log.Panic(err)
			}
		}

		b, err := tx.CreateBucket([]byte(UTXOSetBucketNameYS))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {

			for txIDStr, outs := range utxoMap {
				txID, _ := hex.DecodeString(txIDStr)
				b.Put(txID, gobEncode(outs))
			}
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

/*
	查询对应地址的余额
 */
func (utxoSet *UTXOSetYS) GetBalanceYS(address string) int64 {
	var total int64

	utxos := utxoSet.FindUnspentUTXOsByAddressYS(address)

	for _, utxo := range utxos {
		total += utxo.OutputYS.ValueYS
	}

	return total
}

/*
	查询对应地址, 已打包的UTXO
 */
func (utxoSet *UTXOSetYS) FindUnspentUTXOsByAddressYS(address string) []*UTXOYS {
	var utxos []*UTXOYS

	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	//读数据库
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOSetBucketNameYS))
		if b != nil {
			//遍历UTXOSetBucketName 表
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				//反序列
				txOutputs := DeserializeTxOutputsYS(v)
				//遍历utxos
				for _, utxo := range txOutputs.UTXOsYS {
					//判断地址是否对应
					if utxo.OutputYS.UnlockWithAddressYS(address) {
						utxos = append(utxos, utxo)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return utxos
}

/*
	查询本次转账将要使用的UTXO
 */
func (utxoSet *UTXOSetYS) FindSpentableUTXOsYS(from string, amount int64, txs []*TransactionYS) (int64, map[string][]int) {
	var total int64
	spentableUTXOMap := make(map[string][]int)

	//未打包的UTXO
	unPacketUTXO := utxoSet.FindUnpacketUTXOYS(from, txs)
	for utxo := range  unPacketUTXO {
		fmt.Println("FindUnpacketUTXO")
		fmt.Println(utxo)
	}

	for _, utxo := range unPacketUTXO {
		total += utxo.OutputYS.ValueYS
		txIDStr := hex.EncodeToString(utxo.TxIDYS)
		spentableUTXOMap[txIDStr] = append(spentableUTXOMap[txIDStr], utxo.IndexYS)

		if total >= amount {
			return total, spentableUTXOMap
		}

	}

	//已在区块的UTXO
	packedUTXO := utxoSet.FindUnspentUTXOsByAddressYS(from)

	for _, utxo := range packedUTXO {
		total += utxo.OutputYS.ValueYS
		txIDStr := hex.EncodeToString(utxo.TxIDYS)
		spentableUTXOMap[txIDStr] = append(spentableUTXOMap[txIDStr], utxo.IndexYS)

		if total >= amount {
			return total, spentableUTXOMap
		}
	}

	if total < amount {
		fmt.Printf("%s 的余额不足, 无法转账. 余额为: %d", from, total)
		os.Exit(1)
	}

	return total, spentableUTXOMap
}

/*
	查找对应地址,未打包的UTXO
 */
func (utxoSet *UTXOSetYS) FindUnpacketUTXOYS(from string, txs []*TransactionYS) []*UTXOYS {

	//存储未花费的TxOutput
	var utxos [] *UTXOYS
	//存储已经花费的信息
	spentTxOutputMap := make(map[string][]int) // map[TxID] = []int{vout}

	for i := len(txs) - 1; i >= 0; i-- {
		tx := txs[i]

		utxos = caculateYS(tx, from, spentTxOutputMap, utxos)
	}

	return utxos
}

/*
	更新数据库UTXO
 */
func (utxoSet *UTXOSetYS) UpdateYS() {
	//对最后一个区块进行处理
	lastBlock := utxoSet.blockChainYS.IteratorYS().NextYS()

	//遍历TXs 获取所有input
	txInputs := []*TXInputYS{}
	for _, tx := range lastBlock.TxsYS {
		if !tx.IsCoinBaseTransactionYS() {
			for _, input := range tx.VinsYS {
				txInputs = append(txInputs, input)
			}
		}
	}

	//遍历TXs 获取UTXO
	outsMap := make(map[string]*TxOutputsYS)
	for _, tx := range lastBlock.TxsYS {
		//每个交易中的utxo数组
		utxos := []*UTXOYS{}
		for outIndex, txOut := range tx.VoutsYS {
			isSpent := false
			for _, txInput := range txInputs {
				if txInput.VoutYS == outIndex &&
					bytes.Compare(txInput.TxIDYS, tx.TxIDYS) == 0 {
					//已花费
					isSpent = true
					break
				}
			}
			if isSpent == false {
				utxo := &UTXOYS{tx.TxIDYS, outIndex, txOut}
				utxos = append(utxos, utxo)
			}
		}

		if len(utxos) > 0 {
			txIDStr := hex.EncodeToString(tx.TxIDYS)
			outputs := &TxOutputsYS{utxos}
			outsMap[txIDStr] = outputs
		}
	}

	//for txid, outputs := range outsMap {
	//	fmt.Printf("---------txID :%s", txid)
	//	for _, utxo := range outputs.UTXOs {
	//		fmt.Println(utxo)
	//	}
	//}

	//获取utxo表,将input对应的utxo删除, 添加outsMap中的utxo
	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UTXOSetBucketNameYS))
		if b != nil {
			//1.删除inputs对应的utxo
			for _, input := range txInputs {
				txOutputsBytes := b.Get(input.TxIDYS)
				if len(txOutputsBytes) == 0 {
					continue
				}

				txOutputs := DeserializeTxOutputsYS(txOutputsBytes)

				//是否需要被删除
				isNeedDelete := false

				//将当前outputs里面未被消费的utxo 保存起来
				utxos := []*UTXOYS{}

				for _, utxo := range txOutputs.UTXOsYS {
					if bytes.Compare(utxo.TxIDYS, input.TxIDYS) == 0 &&
						input.VoutYS == utxo.IndexYS &&
						input.UnlockWithAddressYS(utxo.OutputYS.PubKeyHashYS) {
						//已花费
						isNeedDelete = true
						continue
					}

					utxos = append(utxos, utxo)
				}

				if isNeedDelete {
					err := b.Delete(input.TxIDYS)
					if err != nil {
						log.Panic(err)
					}

					if len(utxos) > 0 {
						outputs := &TxOutputsYS{utxos}
						b.Put(input.TxIDYS, gobEncode(outputs))
					}
				}

			}

			//2.添加outsMap 到数据库中
			for txID, outputs := range outsMap {
				txIDBytes, _ := hex.DecodeString(txID)
				b.Put(txIDBytes, gobEncode(outputs))
			}

		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
