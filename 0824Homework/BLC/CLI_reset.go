package BLC

import "os"

func (cli *CLIYS)ResetYS(nodeID string)  {
	blockchain := BlockchainObjectYS(nodeID)
	//defer blockchain.DB.Close()

	if blockchain == nil {
		os.Exit(1)
	}

	utxoSet := &UTXOSetYS{blockchain}
	utxoSet.ResetUTXOSetYS()
}
