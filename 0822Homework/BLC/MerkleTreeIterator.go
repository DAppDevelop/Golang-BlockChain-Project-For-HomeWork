package BLC

type MerkleTreeIteratorYS struct {
	nodes []*MerkleNodeYS
}

func (iterator *MerkleTreeIteratorYS) Next() []*MerkleNodeYS {
	//每次循环新建一个newNodes 保存此层node的数组
	var newNodes []*MerkleNodeYS

	//returnNodes 返回当前层的节点s
	returnNodes := make([]*MerkleNodeYS, len(iterator.nodes), len(iterator.nodes))
	copy(returnNodes, iterator.nodes)

	//fmt.Printf("return: %p\n",returnNodes)
	//fmt.Printf("nodes :%p\n",iterator.nodes)
	//for _, node := range iterator.nodes {
	//	//	fmt.Println("handle node:", node)
	//	//}

	//节点数大于1,则进行合并生成新的节点
	if len(iterator.nodes) > 1 {
		//保证节点为双数
		if len(iterator.nodes)%2 != 0 {
			//fmt.Println("pre nodes count:", len(iterator.nodes))
			iterator.nodes = append(iterator.nodes, iterator.nodes[len(iterator.nodes)-1])
		}
		//fmt.Println("after nodes count:", len(iterator.nodes))
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

	//for _, node := range iterator.nodes {
	//	fmt.Println("return node", node)
	//}

	return returnNodes
}

/*
go run main.go send -from '["12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp","12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp","12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp","12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp"]' -to '["12qg59rrjzA1DFGUTXB5Cbbkh244RtJydwQDhGcvrJpP5VcfwTP","123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ","123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ","123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ"]' -amount '["1","1","1","1"]'
handle node: &{<nil> <nil> [45 94 16 196 52 232 169 247 46 108 255 96 133 41 101 41 246 57 235 161 138 41 46 68 118 163 101 126 104 96 243 22]}
handle node: &{<nil> <nil> [63 175 66 236 240 154 159 108 72 147 148 166 42 12 41 136 53 47 238 136 99 244 51 56 250 7 137 198 61 236 112 46]}
handle node: &{<nil> <nil> [42 66 95 85 174 150 187 202 255 52 128 164 27 195 79 25 98 80 124 244 19 100 20 158 127 150 218 236 241 213 255 27]}
handle node: &{<nil> <nil> [95 146 71 64 7 248 130 79 93 195 201 1 25 249 230 184 149 236 204 22 219 68 72 91 99 182 119 243 214 211 173 217]}
handle node: &{<nil> <nil> [49 216 29 94 8 244 207 36 119 117 24 222 10 47 224 89 195 32 185 224 1 125 195 63 55 232 181 85 56 58 81 110]}
handle node: &{<nil> <nil> [49 216 29 94 8 244 207 36 119 117 24 222 10 47 224 89 195 32 185 224 1 125 195 63 55 232 181 85 56 58 81 110]}
after nodes count: 6
return node &{0xc420186450 0xc420186480 [106 80 132 29 137 104 17 69 127 7 179 253 202 58 28 158 254 50 21 236 61 145 217 221 123 10 75 217 132 179 197 142]}
return node &{0xc4201864b0 0xc4201864e0 [237 17 122 52 171 181 153 60 123 134 3 191 48 0 222 213 218 10 7 37 28 83 9 15 136 94 126 189 97 144 207 56]}
return node &{0xc420186510 0xc420186540 [177 187 117 87 52 25 93 72 149 217 235 161 145 30 52 71 217 140 134 45 125 20 150 115 162 169 173 35 122 146 208 73]}
handle node: &{0xc420186450 0xc420186480 [106 80 132 29 137 104 17 69 127 7 179 253 202 58 28 158 254 50 21 236 61 145 217 221 123 10 75 217 132 179 197 142]}
handle node: &{0xc4201864b0 0xc4201864e0 [237 17 122 52 171 181 153 60 123 134 3 191 48 0 222 213 218 10 7 37 28 83 9 15 136 94 126 189 97 144 207 56]}
handle node: &{0xc420186510 0xc420186540 [177 187 117 87 52 25 93 72 149 217 235 161 145 30 52 71 217 140 134 45 125 20 150 115 162 169 173 35 122 146 208 73]}
pre nodes count: 3
after nodes count: 4
return node &{0xc4201866c0 0xc4201866f0 [83 137 58 106 18 25 215 251 222 108 160 99 47 170 235 42 223 51 171 33 229 108 82 10 93 84 99 110 29 13 254 41]}
return node &{0xc420186720 0xc420186720 [34 60 251 91 147 2 65 140 110 180 60 131 13 220 150 33 197 171 105 122 88 124 200 47 215 147 155 227 116 93 220 92]}
handle node: &{0xc4201866c0 0xc4201866f0 [83 137 58 106 18 25 215 251 222 108 160 99 47 170 235 42 223 51 171 33 229 108 82 10 93 84 99 110 29 13 254 41]}
handle node: &{0xc420186720 0xc420186720 [34 60 251 91 147 2 65 140 110 180 60 131 13 220 150 33 197 171 105 122 88 124 200 47 215 147 155 227 116 93 220 92]}
after nodes count: 2
return node &{0xc420186870 0xc4201868a0 [209 207 185 109 93 50 149 229 168 82 1 213 99 34 237 61 92 214 206 2 9 44 58 110 22 132 216 181 161 35 255 157]}
handle node: &{0xc420186870 0xc4201868a0 [209 207 185 109 93 50 149 229 168 82 1 213 99 34 237 61 92 214 206 2 9 44 58 110 22 132 216 181 161 35 255 157]}
return node &{0xc420186870 0xc4201868a0 [209 207 185 109 93 50 149 229 168 82 1 213 99 34 237 61 92 214 206 2 9 44 58 110 22 132 216 181 161 35 255 157]}
 */
