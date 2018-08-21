package BLC

type UTXOYS struct {
	TxIDYS   []byte      //该output所在的交易的txID
	IndexYS  int         //该output的下标
	OutputYS *TXOutputYS //未花费的output
}
