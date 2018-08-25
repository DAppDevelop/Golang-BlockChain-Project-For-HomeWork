package BLC

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"encoding/gob"
	"bytes"
)

type TransactionPoolYS struct {
	TxsYS []*TransactionYS
}


func NewTXPoolYS(nodeID string) *TransactionPoolYS {
	txPollFile := fmt.Sprintf(txPollFileYS,nodeID)
	//step1：钱包文件不存在
	if _, err := os.Stat(txPollFile); os.IsNotExist(err) {
		fmt.Println("交易池不存在。。。创建交易池")
		txp := &TransactionPoolYS{[]*TransactionYS{}}
		return txp
	}

	txpBytes, err := ioutil.ReadFile(txPollFile)
	if err != nil {
		log.Panic(err)
	}

	var txp TransactionPoolYS

	reader := bytes.NewReader(txpBytes)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&txp)
	if err != nil {

		log.Panic(err)
	}
	return &txp
}


func (txp *TransactionPoolYS) saveFileYS (nodeID string) {
	//组合文件名
	txPollFile := fmt.Sprintf(txPollFileYS,nodeID)
	//将序列化后的ws对象存入文件

	txpBytes := gobEncode(txp)
	err := ioutil.WriteFile(txPollFile, txpBytes, 0644)
	if err != nil {
		log.Panic(err)
	}
}