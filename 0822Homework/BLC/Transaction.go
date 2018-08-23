package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"encoding/hex"
	"math/big"
	"crypto/ecdsa"
	"crypto/rand"
	"time"
	"crypto/elliptic"
)

type TransactionYS struct {
	TxIDYS  []byte        //1. 交易hash
	VinsYS  []*TXInputYS  //2. 输入s
	VoutsYS []*TXOutputYS //3. 输出s
}

/*
	产生创世区块时的Transaction
 */
func NewCoinbaseTransacionYS(address string) *TransactionYS {
	//创建创世区块交易的Vin
	txInput := &TXInputYS{[]byte{}, -1, nil, nil}
	//创建创世区块交易的Vout
	txOutput := NewTxOutputYS(10, address)
	//生产交易Transaction
	txCoinBaseTransaction := &TransactionYS{[]byte{}, []*TXInputYS{txInput}, []*TXOutputYS{txOutput}}
	//设置Transaction的TxHash
	txCoinBaseTransaction.SetIDYS()

	return txCoinBaseTransaction

}


/*
	产生挖矿奖励交易Transaction  奖励为1
 */
func NewRewardTransacionYS(address string) *TransactionYS {
	//创建创世区块交易的Vin
	txInput := &TXInputYS{[]byte{}, -1, nil, nil}
	//创建创世区块交易的Vout
	txOutput := NewTxOutputYS(1, address)
	//生产交易Transaction
	txCoinBaseTransaction := &TransactionYS{[]byte{}, []*TXInputYS{txInput}, []*TXOutputYS{txOutput}}
	//设置Transaction的TxHash
	txCoinBaseTransaction.SetIDYS()

	return txCoinBaseTransaction

}


/*
	创建普通交易Transaction
 */
func NewSimpleTransationYS(from string, to string, amount int64, utxoSet *UTXOSetYS, txs []*TransactionYS) *TransactionYS {
	//1.定义Input和Output的数组
	var txInputs []*TXInputYS
	var txOutputs []*TXOutputYS

	//获取本次转账要使用output
	//total, spentableUTXO := bc.FindSpentableUTXOsYS(from, amount, txs)
	total, spentableUTXO := utxoSet.FindSpentableUTXOsYS(from, amount, txs)

	//2.创建Input
	//获取钱包的集合：
	wallets := NewWalletsYS()
	wallet := wallets.WalletMapYS[from]

	for txID, indexArray := range spentableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TXInputYS{txIDBytes, index, nil, wallet.PublickKeyYS}
			txInputs = append(txInputs, txInput)
		}
	}

	txOutput := NewTxOutputYS(amount, to)
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput2 := NewTxOutputYS(total-amount, from)
	txOutputs = append(txOutputs, txOutput2)

	tx := &TransactionYS{[]byte{}, txInputs, txOutputs}
	tx.SetIDYS()
	//fmt.Println(tx)

	//设置签名
	utxoSet.blockChainYS.SignTrasanctionYS(tx, wallet.PrivateKeyYS, txs)
	return tx
}

func (tx *TransactionYS) IsCoinBaseTransactionYS() bool {
	return len(tx.VinsYS[0].TxIDYS) == 0 && tx.VinsYS[0].VoutYS == -1
}

//签名
/*
签名：为了对一笔交易进行签名
	私钥：
	要获取交易的Input，引用的output，所在的之前的交易：
 */
func (tx *TransactionYS) SignYS(privateKey ecdsa.PrivateKey, prevTxsmap map[string]*TransactionYS) {
	//1.判断当前tx是否时coinbase交易
	if tx.IsCoinBaseTransactionYS() {
		return
	}

	//2.获取input对应的output所在的tx，如果不存在，无法进行签名
	for _, input := range tx.VinsYS {
		if prevTxsmap[hex.EncodeToString(input.TxIDYS)] == nil {
			log.Panic("当前的Input，没有找到对应的output所在的Transaction，无法签名。。")
		}
	}

	//即将进行签名:私钥，要签名的数据
	txCopy := tx.TrimmedCopyYS()

	for index, input := range txCopy.VinsYS {
		// input--->5566

		prevTx := prevTxsmap[hex.EncodeToString(input.TxIDYS)]

		txCopy.VinsYS[index].SignatureYS = nil                                       //仅仅是一个双重保险，保证签名一定为空
		txCopy.VinsYS[index].PublicKeyYS = prevTx.VoutsYS[input.VoutYS].PubKeyHashYS //设置input中的publickey为对应的output的公钥哈希

		txCopy.TxIDYS = txCopy.NewTxIDYS() //产生要签名的数据：

		//为了方便下一个input，将数据再置为空
		txCopy.VinsYS[index].PublicKeyYS = nil

		//获取要交易的数据

		/*
		第一个参数
		第二个参数：私钥
		第三个参数：要签名的数据


		func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error)
		r + s--->sign
		input.Signatrue = sign
	 */
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TxIDYS)
		if err != nil {
			log.Panic(err)
		}

		sign := append(r.Bytes(), s.Bytes()...)
		tx.VinsYS[index].SignatureYS = sign
	}

}

//获取要签名tx的副本
/*
要签名tx中，并不是所有的数据都要作为签名数据，生成签名
txCopy = tx{签名所需要的部分数据}
TxID

Inputs
	txid,vout,sign,publickey

Outputs
	value,pubkeyhash


交易的副本中包含的数据
	包含了原来tx中的输入和输出。
		输入中：sign，publickey
 */

func (tx *TransactionYS) TrimmedCopyYS() *TransactionYS {
	var inputs [] *TXInputYS
	var outputs [] *TXOutputYS
	for _, in := range tx.VinsYS {
		inputs = append(inputs, &TXInputYS{in.TxIDYS, in.VoutYS, nil, nil})
	}

	for _, out := range tx.VoutsYS {
		outputs = append(outputs, &TXOutputYS{out.ValueYS, out.PubKeyHashYS})
	}

	txCopy := &TransactionYS{tx.TxIDYS, inputs, outputs}
	return txCopy

}

/*
	将Transaction 序列化再进行 hash
 */
func (tx *TransactionYS) SetIDYS() {

	txBytes := tx.SerializeYS()

	allBytes := bytes.Join([][]byte{txBytes, IntToHexYS(time.Now().Unix())}, []byte{})

	hash := sha256.Sum256(allBytes)
	//fmt.Printf("transationHash: %x", hash)
	tx.TxIDYS = hash[:]
}

func (tx *TransactionYS) NewTxIDYS() []byte {
	txCopy := tx
	txCopy.TxIDYS = []byte{}
	hash := sha256.Sum256(txCopy.SerializeYS())
	return hash[:]
}

func (tx *TransactionYS) SerializeYS() [] byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

//验证交易
/*
验证的原理：
公钥 + 要签名的数据 验证 签名：rs
 */
func (tx *TransactionYS) VerifityYS(prevTxs map[string]*TransactionYS) bool {
	//1.如果时coinbase交易，不需要验证
	if tx.IsCoinBaseTransactionYS() {
		return true
	}

	//判断当前input是否有对应的Transaction
	for _, input := range tx.VinsYS { //
		if prevTxs[hex.EncodeToString(input.TxIDYS)] == nil {
			log.Panic("当前的input没有找到对应的Transaction，无法验证")
		}
	}

	//验证
	txCopy := tx.TrimmedCopyYS()

	curve := elliptic.P256() //曲线

	for index, input := range tx.VinsYS {
		//原理：再次获取 要签名的数据  + 公钥哈希 + 签名
		/*
		验证签名的有效性：
		第一个参数：公钥
		第二个参数：签名的数据
		第三、四个参数：签名：r，s
		func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool
		 */
		//ecdsa.Verify()

		//获取要签名的数据
		prevTx := prevTxs[hex.EncodeToString(input.TxIDYS)]

		txCopy.VinsYS[index].SignatureYS = nil
		txCopy.VinsYS[index].PublicKeyYS = prevTx.VoutsYS[input.VoutYS].PubKeyHashYS
		txCopy.TxIDYS = txCopy.NewTxIDYS() //要签名的数据

		txCopy.VinsYS[index].PublicKeyYS = nil

		//获取公钥
		/*
		type PublicKey struct {
			elliptic.Curve
			X, Y *big.Int
		}
		 */

		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PublicKeyYS)
		x.SetBytes(input.PublicKeyYS[:keyLen/2])
		y.SetBytes(input.PublicKeyYS[keyLen/2:])

		rawPublicKey := ecdsa.PublicKey{curve, &x, &y}

		//获取签名：

		r := big.Int{}
		s := big.Int{}

		signLen := len(input.SignatureYS)
		r.SetBytes(input.SignatureYS[:signLen/2])
		s.SetBytes(input.SignatureYS[signLen/2:])

		if ecdsa.Verify(&rawPublicKey, txCopy.TxIDYS, &r, &s) == false {
			return false
		}

	}
	return true
}

//格式化输出
func (tx *TransactionYS) String() string {
	var vinStrings [][]byte
	for _, vin := range tx.VinsYS {
		vinString := fmt.Sprint(vin)
		vinStrings = append(vinStrings, []byte(vinString))
	}
	vinString := bytes.Join(vinStrings, []byte{})

	var outStrings [][]byte
	for _, out := range tx.VoutsYS {
		outString := fmt.Sprint(out)
		outStrings = append(outStrings, []byte(outString))
	}

	outString := bytes.Join(outStrings, []byte{})

	return fmt.Sprintf("\n\r\t\t===============================\n\r\t\tTxID: %x, \n\t\tVins: %v, \n\t\tVout: %v\n\t\t", tx.TxIDYS, string(vinString), string(outString))
}
