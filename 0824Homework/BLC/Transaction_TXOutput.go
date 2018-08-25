package BLC

import (
	"fmt"
	"bytes"
)

type TXOutputYS struct {
	ValueYS        int64  //金额
	//ScriptPubKey string //用户名(scriptPubkey:锁定脚本,包含公钥)
	PubKeyHashYS [] byte//公钥哈希
}

//判断TxOutput是否时指定的用户解锁
func (txOutput *TXOutputYS) UnlockWithAddressYS(address string) bool{
	full_payload:=Base58Decode([]byte(address))

	pubKeyHash:=full_payload[1:len(full_payload)-addressCheckSumLenYS]

	return bytes.Compare(pubKeyHash,txOutput.PubKeyHashYS) == 0
}

//根据地址创建一个output对象
func NewTxOutputYS(value int64,address string) *TXOutputYS{
	txOutput:=&TXOutputYS{value,nil}
	txOutput.LockYS(address)
	return txOutput
}

//锁定
func (tx *TXOutputYS) LockYS(address string){
	full_payload := Base58Decode([]byte(address))
	//获取公钥hash
	tx.PubKeyHashYS = full_payload[1:len(full_payload)-addressCheckSumLenYS]
}

//格式化输出
func (tx *TXOutputYS) String() string {
	return fmt.Sprintf("\n\t\t\tValue: %d, PubKeyHash(转成地址显示): %s", tx.ValueYS, PublicHashToAddressYS(tx.PubKeyHashYS))
	//return fmt.Sprintf("\n\t\t\tValue: %d, PubKeyHash: %x", tx.Value, tx.PubKeyHash)
}
