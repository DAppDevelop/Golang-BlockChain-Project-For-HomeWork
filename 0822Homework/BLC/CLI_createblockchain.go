package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIYS) createGenesisBlockchainYS(address string) {
	if IsValidAddressYS([]byte(address)) {
		CreateBlockchainWithGenesisBlockYS(address)

		block := BlockchainObjectYS()
		defer block.DBYS.Close()

		if block == nil {
			fmt.Println("没有数据库。。")
			os.Exit(1)
		}

		utxoSet := &UTXOSetYS{block}
		utxoSet.ResetUTXOSetYS()

	} else {
		fmt.Println("地址无效")

	}



}