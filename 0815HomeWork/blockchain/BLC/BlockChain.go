package BLC

import "log"

type BlockchainYS struct {
	BlockYSs []*BlockYS
}

//1. 创建带有创世区块的区块链
func CreatBlockchainWithGenesisBlockYS() *BlockchainYS {
	// 创建创世区块
	genesisBlock := CreateGenesisBlockYS("Genesis Data....... ")
	// 返回区块链对象
	return &BlockchainYS{[]*BlockYS{genesisBlock}}
}

// 增加区块到区块链里面
func (blc *BlockchainYS) AddBlockToBlockchainYS(data string, height int64, preHash []byte) {
	// 创建新区块
	newBlock := NewBlockYS(data, height, preHash)

	pof:= NewProofOfWorkYS(newBlock)
	//判断工作量证明是否有效
	if pof.IsValidYS() {
		// 往链里面添加区块
		blc.BlockYSs = append(blc.BlockYSs, newBlock)
	} else {
		log.Panic("工作量证明无效")
	}
}
