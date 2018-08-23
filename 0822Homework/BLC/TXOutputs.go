package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOutputsYS struct {
	UTXOsYS []*UTXOYS
}

func (outs *TxOutputsYS) SerializeYS() []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeTxOutputsYS(data []byte) *TxOutputsYS  {
	outs := TxOutputsYS{}

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&outs)

	if err != nil {
		log.Panic(err)
	}

	return &outs
}