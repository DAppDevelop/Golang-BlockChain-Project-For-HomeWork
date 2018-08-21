package BLC

import "fmt"

type TXInputYS struct {
	TxIDYS      []byte //交易的Hash
	VoutYS      int    //存储TXOutput在Vout里面的索引(第几个交易)
	ScriptSigYS string //用户名花费的是谁的钱(解锁脚本,包含数字签名)
}

//判断TXInput是否指定的address消费
func (txInput *TXInputYS) UnlockWithAddressYS(address string) bool {
	return txInput.ScriptSigYS == address
}

//格式化输出
func (tx *TXInputYS) String() string {
	return fmt.Sprintf("\n\t\t\tTxInput_TXID: %x, Vout: %v, ScriptSig: %v", tx.TxIDYS, tx.VoutYS, tx.ScriptSigYS)
}
