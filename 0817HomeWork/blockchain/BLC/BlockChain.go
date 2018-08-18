package BLC

import (
	"log"
	"os"
	"github.com/boltdb/bolt"
	"fmt"
	"math/big"
)

type BlockchainYS struct {
	TipYS []byte //最新的区块的Hash
	DBYS  *bolt.DB
}

//1. 创建带有创世区块的区块链
func CreatBlockchainWithGenesisBlockYS(data string) {
	//判断数据库是否已经存
	if DBExistsYS() {
		fmt.Println("Genesis Block 已经存在...")
		os.Exit(1)
	}

	fmt.Println("创建创世区块....")

	//创建或打开数据库
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {

		//创建表
		b, err := tx.CreateBucket([]byte(BlockBucketName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建创世区块
			genesisBlock := CreateGenesisBlockYS(data)

			//序列号block并存入数据库
			err := b.Put([]byte(genesisBlock.HashYS), []byte(genesisBlock.SerializeYS()))

			if err != nil {
				log.Panic(err)
			}

			//更新数据库最新区块hash
			err = b.Put([]byte("l"), []byte(genesisBlock.HashYS))

			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

// 增加区块到区块链里面
func (blc *BlockchainYS) AddBlockToBlockchainYS(data string) {
	err := blc.DBYS.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))

		if b != nil {
			//取到最新区块
			blockbyte := b.Get(blc.TipYS)

			block := DeserializeBlockYS(blockbyte)

			// 创建新区块
			newBlock := NewBlockYS(data, block.HeightYS+1, block.HashYS)

			//序列号block并存入数据库
			err := b.Put(newBlock.HashYS, newBlock.SerializeYS())

			if err != nil {
				log.Panic(err)
			}

			//更新数据库最新区块hash
			err = b.Put([]byte("l"), []byte(newBlock.HashYS))

			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

// 遍历输出所有区块的信息
func (blc *BlockchainYS) PrintchainYS() {
	//创建迭代器
	blockIterator := blc.IteratorYS()

	for {
		block := blockIterator.NextYS()

		fmt.Println(block)

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHashYS)

		//判断当期的block是否为创世区块（创世区块perblockhash为000000....）
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

//迭代器
func (blockchain *BlockchainYS) IteratorYS() *BlockchainIteratorYS {
	return &BlockchainIteratorYS{blockchain.TipYS, blockchain.DBYS}
}

// 判断数据库是否存在
func DBExistsYS() bool {
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		return false
	}

	return true
}

// 返回Blockchain对象
func BlockchainObjectYS() *BlockchainYS {
	//因为已经知道数据库的名字，所以只要取出最新区块hash，既可以返回blockchain对象
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockBucketName))
		if b != nil {
			//取出最新区块hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &BlockchainYS{tip, db}
}
