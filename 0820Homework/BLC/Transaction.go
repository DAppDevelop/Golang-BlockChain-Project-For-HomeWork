package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
)

type TransactionYS struct {
	TxIDYS  []byte        //1. 交易hash
	VinsYS  []*TXInputYS  //2. 输入s
	VoutsYS []*TXOutputYS //3. 输出s
}

/*
	产生创世区块时的Transaction
 */
func NewCoinbaseTransacionYS(address string) *TransactionYS {
	//创建创世区块交易的Vin
	txInput := &TXInputYS{[]byte{}, -1, "Genesis DATA"}
	//创建创世区块交易的Vout
	txOutput := &TXOutputYS{10, address}
	//生产交易Transaction
	txCoinBaseTransaction := &TransactionYS{[]byte{}, []*TXInputYS{txInput}, []*TXOutputYS{txOutput}}
	//设置Transaction的TxHash
	txCoinBaseTransaction.SetIDYS()

	return txCoinBaseTransaction

}

/*
	创建普通交易Transaction
 */
func NewSimpleTransationYS(from string, to string, amount int64, bc *BlockchainYS, txs []*TransactionYS) *TransactionYS {
	//1.定义Input和Output的数组
	var txInputs []*TXInputYS
	var txOutputs []*TXOutputYS

	//获取本次转账要使用output
	total, spentableUTXO := bc.FindSpentableUTXOsYS(from, amount, txs)

	//2.创建Input
	for txID, indexArray := range spentableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TXInputYS{txIDBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutput := &TXOutputYS{amount, to}
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput2 := &TXOutputYS{total - amount, from}
	txOutputs = append(txOutputs, txOutput2)

	tx := &TransactionYS{[]byte{}, txInputs, txOutputs}
	tx.SetIDYS()
	//fmt.Println(tx)
	return tx
}

/*
	将Transaction 序列化再进行 hash
 */
func (tx *TransactionYS) SetIDYS() {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash := sha256.Sum256(result.Bytes())
	//fmt.Printf("transationHash: %x", hash)
	tx.TxIDYS = hash[:]
}

func (tx *TransactionYS) IsCoinBaseTransactionYS() bool {
	return len(tx.VinsYS[0].TxIDYS) == 0 && tx.VinsYS[0].VoutYS == -1
}

//格式化输出
func (tx *TransactionYS) String() string {
	var vinStrings [][]byte
	for _, vin := range tx.VinsYS {
		vinString := fmt.Sprint(vin)
		vinStrings = append(vinStrings, []byte(vinString))
	}
	vinString := bytes.Join(vinStrings, []byte{})

	var outStrings [][]byte
	for _, out := range tx.VoutsYS {
		outString := fmt.Sprint(out)
		outStrings = append(outStrings, []byte(outString))
	}

	outString := bytes.Join(outStrings, []byte{})

	return fmt.Sprintf("\n\r\t\t===============================\n\r\t\tTxID: %x, \n\t\tVins: %v, \n\t\tVout: %v\n\t\t", tx.TxIDYS, string(vinString), string(outString))
}
