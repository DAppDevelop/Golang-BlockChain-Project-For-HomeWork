package BLC

import "fmt"

func (cli *CLIYS) createGenesisBlockchainYS(address string) {
	if IsValidAddressYS([]byte(address)) {
		CreateBlockchainWithGenesisBlockYS(address)
	} else {
		fmt.Println("地址无效")

	}

}