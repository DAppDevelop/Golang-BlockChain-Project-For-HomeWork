package BLC

type InvYS struct {
	AddrFromYS string   //当前节点地址
	TypeYS     string   //类型（block or Transaction
	ItemsYS    [][]byte //对应类型的数据的hash
}
