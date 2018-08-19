package BLC

import "fmt"

type TXOutputYS struct {
	ValueYS        int64  //金额
	ScriptPubKeyYS string //用户名(scriptPubkey:锁定脚本,包含公钥)
}

//判断TXOutput是否指定的address解锁
func (txOutput *TXOutputYS) UnlockWithAddressYS(address string) bool {
	return txOutput.ScriptPubKeyYS == address
}

//格式化输出
func (tx *TXOutputYS) String() string {
	return fmt.Sprintf("\n\t\t\tValue: %d, ScriptPubKey: %s", tx.ValueYS, tx.ScriptPubKeyYS)
}
