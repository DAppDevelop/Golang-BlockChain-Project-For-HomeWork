package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

type ProofOfWorkYS struct {
	Block *BlockYS // 当前要验证的区块

	target *big.Int // 当hash小于此target时，为挖矿成功
}

func NewProofOfWorkYS(block *BlockYS) *ProofOfWorkYS {
	//创建一个初始值为1的target
	target := big.NewInt(1)

	//左移 256 - targetBit
	target = target.Lsh(target, 256-targetBitYS)

	return &ProofOfWorkYS{block, target}
}

func (pow *ProofOfWorkYS) RunYS() ([]byte, int64) {
	//使用nonce计算hash不符合target时候，加1，直到hash符合要求
	nonce := 0

	var hashInt big.Int
	var hash [32]byte

	dataBytes := pow.prepareDataYS() //计算除了nonce以外的block参数
	//fmt.Printf("dataBytes:%x",dataBytes)
	for {
		//将Block的属性拼接成字节数组作为sha256.Sum256的入参
		dataBytes := bytes.Join(
			[][]byte{ //[]byte的切片
				dataBytes,
				IntToHexYS(int64(nonce)),
			},
			[]byte{},
		)
		//生成hash
		hash = sha256.Sum256(dataBytes)
		//fmt.Printf("\r%x", hash)

		//将hash转换成*int类型并返回给hashInt
		hashInt.SetBytes(hash[:])
		//判断hash有效性，如果满足条件，跳出循环
		if pow.target.Cmp(&hashInt) == 1 {
			fmt.Printf("\nhash: %x\n", hash) //hash: 00ea9e3743900b6086acbb86390457f72fb3a4908609bd900536064f8e89448d
			break
		}

		//如果不满足条件，nonce+1并继续循环
		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}

// 数据拼接，返回字节数组
func (pow *ProofOfWorkYS) prepareDataYS() []byte {
	//bytes.Join 以sep为连接符，拼接[][]byte
	data := bytes.Join(
		[][]byte{ //[]byte的切片
			pow.Block.PrevBlockHashYS,
			pow.Block.HashTransactionsYS(),
			IntToHexYS(pow.Block.TimestampYS),
			IntToHexYS(int64(targetBitYS)),
			IntToHexYS(int64(pow.Block.HeightYS)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWorkYS) IsValidYS() bool {

	var hashInt big.Int

	hashInt.SetBytes(pow.Block.HashYS)

	//1.proofOfWork.Block.Hash
	//2.proofOfWork.Target 作比较
	if pow.target.Cmp(&hashInt) == 1 {
		return true
	}

	return false
}
