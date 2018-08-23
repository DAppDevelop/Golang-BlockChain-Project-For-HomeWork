package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIYS) printchainYS() {
	if DBExistsYS() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObjectYS()

	defer blockchain.DBYS.Close()

	blockchain.PrintchainYS()
}
