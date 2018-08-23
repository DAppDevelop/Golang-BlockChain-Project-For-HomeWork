package BLC

type MerkleTreeIteratorYS struct {
	nodes []*MerkleNodeYS
}

func (iterator *MerkleTreeIteratorYS)Next() []*MerkleNodeYS  {
	//每次循环新建一个newNodes 保存此层node的数组
	var newNodes []*MerkleNodeYS

	//returnNodes 返回当前层的节点s
	returnNodes := make([]*MerkleNodeYS, len(iterator.nodes), len(iterator.nodes))
	copy(returnNodes, iterator.nodes)

	//fmt.Printf("return: %p\n",returnNodes)
	//fmt.Printf("nodes :%p\n",iterator.nodes)

	//节点数大于1,则进行合并生成新的节点
	if len(iterator.nodes) > 1 {
		for i := 0; i < len(iterator.nodes); i += 2 {

			node := NewMerkleNode(iterator.nodes[i], iterator.nodes[i+1], nil)
			//fmt.Printf("node: %x = %x + %x", node.DataHashYS, iterator.nodes[i].DataHashYS, iterator.nodes[i+1].DataHashYS)
			newNodes = append(newNodes, node)

		}
		//设置新的nodes
		iterator.nodes = newNodes
	}

	//fmt.Printf("return: %p\n",returnNodes)
	//fmt.Printf("nodes :%p\n",iterator.nodes)

	//for _, node := range newNodes {
	//	fmt.Println(node)
	//}

	return returnNodes
}
