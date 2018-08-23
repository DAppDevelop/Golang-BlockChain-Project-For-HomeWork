package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type BlockYS struct {
	HeightYS        int64            //1. 区块高度
	PrevBlockHashYS []byte           //2. 上一个区块HASH
	TxsYS           []*TransactionYS //3. 交易数据
	TimestampYS     int64            //4. 时间戳
	HashYS          []byte           //5. Hash
	NonceYS         int64            //6. Nonce
}

func NewBlockYS(txs []*TransactionYS, height int64, preBlockHash []byte) *BlockYS {
	block := &BlockYS{height, preBlockHash, txs, time.Now().Unix(), nil, 0}

	//创建工作量证明结构体
	pow := NewProofOfWorkYS(block)

	//调用工作量证明的方法并且返回有效的Hash和Nonce（挖矿）
	hash, nonce := pow.RunYS()
	block.HashYS = hash[:]
	block.NonceYS = nonce

	return block
}

// 创建创世区块
func CreateGenesisBlockYS(txs []*TransactionYS) *BlockYS {
	return NewBlockYS(txs, 1, make([]byte, 32, 32))
}

// 需要将Txs转换成[]byte(256)
func (block *BlockYS) HashTransactionsYS() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.TxsYS {
		txHashes = append(txHashes, tx.TxIDYS)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]

}

//打印格式
func (block *BlockYS) String() string {
	return fmt.Sprintf(
		"\n------------------------------"+
			"\nABlock's Info:\n\t"+
			"Height:%d,\n\t"+
			"PreHash:%x,\n\t"+
			"Txs: %v,\n\t"+
			"Timestamp: %s,\n\t"+
			"Hash: %x,\n\t"+
			"Nonce: %v\n\t",
		block.HeightYS,
		block.PrevBlockHashYS,
		block.TxsYS,
		time.Unix(block.TimestampYS, 0).Format("2006-01-02 03:04:05 PM"),
		block.HashYS, block.NonceYS)
}

// 序列化：将区块序列化成字节数组
func (block *BlockYS) SerializeYS() []byte {

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	//fmt.Println(result.Bytes())
	return result.Bytes()
}

// 反序列化：将字节数组反序列化为block对象
func DeserializeBlockYS(blockBytes []byte) *BlockYS {

	var block BlockYS

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
