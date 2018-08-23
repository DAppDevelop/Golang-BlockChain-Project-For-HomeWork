package BLC

import (
	"fmt"
	"bytes"
)

type TXOutputYS struct {
	ValueYS    int64   //金额
	PubKeyHashYS [] byte //公钥哈希
	//ScriptPubKeyYS string  //用户名(scriptPubkey:锁定脚本,包含公钥)
}

//根据地址创建一个output对象
func NewTxOutputYS(value int64,address string) *TXOutputYS{
	txOutput:=&TXOutputYS{value,nil}
	txOutput.LockYS(address)
	return txOutput
}

//判断TXOutput是否指定的address解锁(比较地址计算出的哈希公钥是否和output保存PubKeyHash的一致
func (txOutput *TXOutputYS) UnlockWithAddressYS(address string) bool {
	full_payload:=Base58Decode([]byte(address))

	pubKeyHash:=full_payload[1:len(full_payload)-addressCheckSumLenYS]

	return bytes.Compare(pubKeyHash,txOutput.PubKeyHashYS) == 0
}

//锁定(通过地址address获取哈希公钥PubKeyHash)
func (tx *TXOutputYS) LockYS(address string){
	full_payload := Base58Decode([]byte(address))
	//获取公钥hash
	tx.PubKeyHashYS = full_payload[1:len(full_payload)-addressCheckSumLenYS]
}


//格式化输出
func (tx *TXOutputYS) String() string {
	return fmt.Sprintf("\n\t\t\tValue: %d, PubKeyHash(转成地址显示): %s", tx.ValueYS, PublicHashToAddress(tx.PubKeyHashYS))
}
