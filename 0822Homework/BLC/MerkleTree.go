package BLC

import (
	"crypto/sha256"
)

type MerkleNodeYS struct {
	LeftNodeYS  *MerkleNodeYS
	RightNodeYS *MerkleNodeYS
	DataHashYS  []byte
}

type MerkleTreeYS struct {
	RootNodeYS *MerkleNodeYS
}

func NewMerkleNode(lelfNode, rightNode *MerkleNodeYS, txHash []byte) *MerkleNodeYS {
	node := &MerkleNodeYS{}

	var hash [32]byte
	if lelfNode == nil && rightNode == nil {
		//如果是叶子节点
		hash = sha256.Sum256(txHash)
	} else {
		//子节点
		//拼接dataHash
		prevHash := append(lelfNode.DataHashYS, rightNode.DataHashYS...)
		hash = sha256.Sum256(prevHash)
	}

	node.LeftNodeYS = lelfNode
	node.RightNodeYS = rightNode
	node.DataHashYS = hash[:]

	return node

}

func NewMerkleTreeYS(txHashDataYS [][]byte) *MerkleTreeYS {
	//保存每层merkle Tree 节点 当节点数为1时, 跳出循环
	var nodes []*MerkleNodeYS

	merkleTree := &MerkleTreeYS{}

	//创建叶子节点
	//当txHashData 为奇数 , 最后一个复制补全
	if len(txHashDataYS)%2 != 0 {
		txHashDataYS = append(txHashDataYS, txHashDataYS[len(txHashDataYS)-1])
	}

	for _, txHash := range txHashDataYS {
		node := NewMerkleNode(nil, nil, txHash)
		nodes = append(nodes, node)
	}

	iterator := merkleTree.Iterator(nodes)

	//生成子节点(循环到根节点生成为止)
	for {
		newNodes := iterator.Next()

		//判断当前层node的数量是否为1, 为1则为根节点
		if len(newNodes) == 1 {
			merkleTree.RootNodeYS = newNodes[0]
			//fmt.Printf("newNodes[0].DataHashYS: %x\n", merkleTree.RootNodeYS.DataHashYS)
			return merkleTree
		}
	}

	return &MerkleTreeYS{}

}

func (mt *MerkleTreeYS)Iterator(nodes []*MerkleNodeYS) *MerkleTreeIteratorYS {
	return &MerkleTreeIteratorYS{nodes}
}
