package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"os"
)

type BlockchainIteratorYS struct {
	currentHashYS []byte   //当前hash
	DBYS          *bolt.DB //数据库
}

/*
	根据当前迭代器currentHash从数据库中查找对应的block,
	之后将迭代器的currentHash置为前一个区块hash.
 */
func (blockchainIterator *BlockchainIteratorYS) NextYS() *BlockYS {
	var block BlockYS
	DBName := fmt.Sprintf(DBNameYS, os.Getenv("NODE_ID"))
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketNameYS))
		if b != nil {
			//获取当期迭代器对应的block
			currentBlockBytes := b.Get(blockchainIterator.currentHashYS)
			//block = DeserializeBlock(currentBlockBytes)
			gobDecode(currentBlockBytes, &block)

			//将迭代器的currentHash 置为 上一个区块的hash
			blockchainIterator.currentHashYS = block.PrevBlockHashYS
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &block
}
