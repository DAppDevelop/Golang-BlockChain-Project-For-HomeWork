package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockchainIteratorYS struct {
	currentHashYS []byte
	DBYS          *bolt.DB
}

func (blockchainIterator *BlockchainIteratorYS) NextYS() *BlockYS {
	var block *BlockYS

	err := blockchainIterator.DBYS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {
			//获取当期迭代器对应的block
			currentBlockBytes := b.Get(blockchainIterator.currentHashYS)
			block = DeserializeBlockYS(currentBlockBytes)

			//将迭代器的currentHash 置为 上一个区块的hash
			blockchainIterator.currentHashYS = block.PrevBlockHashYS
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return block
}
