package BLC

import (
	"os"
)

func (cli *CLIYS) createGenesisBlockchainYS(address string, nodeID string) {
	CreateBlockchainWithGenesisBlockYS(address,nodeID)

	blockchain := BlockchainObjectYS(nodeID)
	//defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	utxoSet := &UTXOSetYS{blockchain}
	utxoSet.ResetUTXOSetYS()
}