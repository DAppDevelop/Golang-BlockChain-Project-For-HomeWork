package main

import "chenyuanshan/0817Homework/blockchain/BLC"

func main() {

	//blockchain := BLC.CreatBlockchainWithGenesisBlockYS()
	//
	//blockchain.AddBlockToBlockchainYS("2", blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HeightYS+1, blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HashYS )
	//blockchain.AddBlockToBlockchainYS("33", blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HeightYS+1, blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HashYS )
	//blockchain.AddBlockToBlockchainYS("444", blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HeightYS+1, blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HashYS )
	//blockchain.AddBlockToBlockchainYS("5555", blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HeightYS+1, blockchain.BlockYSs[len(blockchain.BlockYSs)-1].HashYS )
	//
	//fmt.Println(blockchain.BlockYSs)

	cli := BLC.CLIYS{}
	cli.RunYS()

}
