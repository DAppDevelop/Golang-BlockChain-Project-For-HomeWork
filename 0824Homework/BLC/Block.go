package BLC

import (
	"time"
	"fmt"
)

type BlockYS struct {
	HeightYS        int64          //1. 区块高度
	PrevBlockHashYS []byte         //2. 上一个区块HASH
	TxsYS           []*TransactionYS //3. 交易数据
	TimestampYS     int64          //4. 时间戳
	HashYS          []byte         //5. Hash
	NonceYS        int64          //6. Nonce
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
	//make([]byte,32,32) --> 32*8=256
	return NewBlockYS(txs, 1, make([]byte, 32, 32))
}

// 需要将Txs里面每个tx.TxID拼接后hash, 转换成[]byte(256)
func (block *BlockYS) HashTransactionsYS() []byte {
	//将txs的hash序列化为[]byte,并放进一个数组里面
	var txs [][]byte
	for _,tx := range block.TxsYS {
		txBytes := gobEncode(tx)
		txs = append(txs, txBytes)
	}

	merkleTree := NewMerkleTreeYS(txs)

	return merkleTree.RootNodeYS.DataHashYS
}


//打印格式
func (block *BlockYS) String() string {
	return fmt.Sprintf(
		"\n------------------------------"+
			"\nBlock's Info:\n\t"+
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
