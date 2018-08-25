package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIYS) getBalanceYS(address string,nodeID string) {
	blockchain := BlockchainObjectYS(nodeID)
	//defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	//txs 传nil值，查询时没有新的交易产生
	//total := blockchain.GetBalance(address, []*Transaction{})
	utxoSet := &UTXOSetYS{blockchain}
	total := utxoSet.GetBalanceYS(address)
	fmt.Printf("%s的余额：%d\n", address, total)
}
