package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIYS) getBalanceYS(address string) {
	if DBExistsYS() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObjectYS()

	defer blockchain.DBYS.Close()

	if !IsValidAddressYS([]byte(address)) {
		fmt.Println("地址无效")
		os.Exit(1)
	}

	//txs 传nil值，查询时没有新的交易产生
	total := blockchain.GetBalanceYS(address, []*TransactionYS{})
	fmt.Printf("%s的余额：%d\n", address, total)
}
