package BLC

import "os"

func (cli *CLIYS) printchainYS(nodeID string) {
	blockchain := BlockchainObjectYS(nodeID)
	//defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	blockchain.PrintchainYS()
}
