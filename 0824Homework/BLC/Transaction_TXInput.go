package BLC

import (
	"fmt"
	"bytes"
)

type TXInputYS struct {
	TxIDYS    []byte // 1. 交易的Hash
	VoutYS      int    //2. 存储TXOutput在Vout里面的索引(第几个交易)
	//ScriptSig string // 3. 用户名花费的是谁的钱(解锁脚本,包含数字签名)
	SignatureYS []byte //数字签名
	PublicKeyYS[]byte //原始公钥，钱包里的公钥
}



//判断TXInput是否指定的address消费
func (txInput *TXInputYS) UnlockWithAddressYS(pubKeyHash []byte) bool {
	pubKeyHash2:=PubKeyHashYS(txInput.PublicKeyYS)
	return bytes.Compare(pubKeyHash,pubKeyHash2) == 0
}

//格式化输出
func (tx *TXInputYS) String() string {
	return fmt.Sprintf("\n\t\t\tTxInput_TXID: %x, Vout: %v, Signature: %x, PublicKey:%x", tx.TxIDYS, tx.VoutYS, tx.SignatureYS, tx.PublicKeyYS)
}