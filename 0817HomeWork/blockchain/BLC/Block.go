package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
)

type BlockYS struct {
	HeightYS        int64  //1. 区块高度
	PrevBlockHashYS []byte //2. 上一个区块HASH
	DataYS          []byte //3. 交易数据
	TimestampYS     int64  //4. 时间戳
	HashYS          []byte //5. Hash
	NonceYS         int64  //6. Nonce
}

func NewBlockYS(data string, height int64, preBlockHash []byte) *BlockYS {
	block := &BlockYS{
		height,
		preBlockHash,
		[]byte(data),
		time.Now().Unix(),
		nil,
		0}

	// 调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWorkYS(block)

	// 挖矿验证
	hash, nonce := pow.RunYS()
	block.HashYS = hash[:]
	block.NonceYS = nonce

	return block

}

func CreateGenesisBlockYS(data string) *BlockYS {
	return NewBlockYS(data, 1, make([]byte, 32, 32))
}

//格式化
func (block *BlockYS) String() string {
	return fmt.Sprintf(
		"\n------------------------------"+
			"\nABlock's Info:\n\t"+
			"Height:%d,\n\t"+
			"PreHash:%x,\n\t"+
			"Data: %s,\n\t"+
			"Timestamp: %s,\n\t"+
			"Hash: %x,\n\t"+
			"Nonce: %v\n\t",
		block.HeightYS, block.PrevBlockHashYS, block.DataYS, time.Unix(block.TimestampYS, 0).Format("2006-01-02 03:04:05 PM"), block.HashYS, block.NonceYS)
}

// 将区块序列化成字节数组
func (block *BlockYS) SerializeYS() []byte {

	var buff bytes.Buffer
	//*buff才是io.Writer接口实现
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// 反序列化
func DeserializeBlockYS(blockBytes []byte) *BlockYS {

	var block BlockYS

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
