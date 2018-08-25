package BLC

type GetDataYS struct {
	AddrFromYS string //当前节点自己的地址
	TypeYS string //数据类型（block或者tx)
	HashYS []byte//block或者Tx的hash
}
