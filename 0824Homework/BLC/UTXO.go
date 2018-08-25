package BLC

import (
	"fmt"
	"encoding/hex"
)

type UTXOYS struct {
	TxIDYS   []byte    //1.该Output所在的交易id
	IndexYS  int       //2.该Output 的下标
	OutputYS *TXOutputYS //3.Output
}


//打印格式
func (utxo *UTXOYS) String() string {
	return fmt.Sprintf(
		"\n------------------------------"+
			"\nA UTXO's Info:\n\t"+
			"TxID:%s,\n\t"+
			"Index:%d,\n\t"+
			"Output: %v,\n\t",
		hex.EncodeToString(utxo.TxIDYS),
		utxo.IndexYS,
		utxo.OutputYS,
		)
}
