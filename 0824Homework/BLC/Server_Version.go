package BLC

//定义为12字节长度
type VersionYS struct {
	VersionYS    int64  //版本
	BestHeightYS int64  //当前节点区块链中最后一个区块的高度
	AddrFromYS   string //当前节点地址
}
