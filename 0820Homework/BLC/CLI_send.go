package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIYS) sendYS(from []string, to []string, amount []string) {
	//go run main.go send -from '["yancey"]' -to '["a"]' -amount '["10"]'
	if DBExistsYS() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObjectYS()
	defer blockchain.DBYS.Close()

	blockchain.MineNewBlockYS(from, to, amount)
}