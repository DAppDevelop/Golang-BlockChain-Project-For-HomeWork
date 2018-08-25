package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIYS)startNodeYS (nodeID string, miner string)  {
	//fmt.Println(miner)
	if miner == "" || IsValidAddressYS([]byte(miner)) {
		startServerYS(nodeID, miner)
	} else {
		fmt.Println("Miner地址无效")
		os.Exit(1)
	}
}
