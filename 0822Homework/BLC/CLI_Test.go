package BLC

func (cli *CLIYS)Test()  {
	bc := BlockchainObjectYS()
	defer bc.DBYS.Close()

	utxoSet := &UTXOSetYS{bc}
	utxoSet.ResetUTXOSetYS()


}
