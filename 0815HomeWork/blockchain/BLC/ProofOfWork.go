package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

type ProofOfWorkYS struct {
	BlockYS  *BlockYS   // 当前要验证的区块
	targetYS *big.Int // 大数据存储 2^24
}

func NewProofOfWorkYS(block *BlockYS) *ProofOfWorkYS {
	//1. 创建一个初始值为1的target
	target := big.NewInt(1)

	//2. 左移256 - targetBit
	target = target.Lsh(target, 256-targetBitYS)

	return &ProofOfWorkYS{block, target}
}

func (proofOfWork *ProofOfWorkYS) RunYS() ([]byte, int64) {
	nonce := 0

	var hashInt big.Int
	var hash [32]byte

	for {
		//1. 将Block的属性拼接成字节数组
		dataBytes := proofOfWork.prepareDataYS(nonce)

		//2. 生成hash
		hash = sha256.Sum256(dataBytes)

		hashInt.SetBytes(hash[:])

		//判断hashInt是否小于Block里面的target
		//3. 判断hash有效性，如果满足条件，跳出循环
		if proofOfWork.targetYS.Cmp(&hashInt) == 1 {
			fmt.Printf("hash: %x\n", hash)
			break
		}

		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}



// 数据拼接，返回字节数组
func (pow *ProofOfWorkYS) prepareDataYS(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.BlockYS.PrevBlockHashYS,
			pow.BlockYS.DataYS,
			IntToHexYS(pow.BlockYS.TimestampYS),
			IntToHexYS(int64(targetBitYS)),
			IntToHexYS(int64(nonce)),
			IntToHexYS(int64(pow.BlockYS.HeightYS)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *ProofOfWorkYS) IsValidYS() bool {

	var hashInt big.Int
	//获取block的hash值并转换为big.Int
	hashInt.SetBytes(proofOfWork.BlockYS.HashYS)
	//比较target，判断是否小于target
	if proofOfWork.targetYS.Cmp(&hashInt) == 1 {
		return true
	}

	return false
}
