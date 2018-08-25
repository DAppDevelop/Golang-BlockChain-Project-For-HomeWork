package BLC

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

type WalletsYS struct {
	WalletMapYS map[string]*WalletYS
}



//提供一个函数，用于创建一个钱包的集合
/*
思路：修改该方法：
	读取本地的钱包文件，如果文件存在，直接获取
	如果文件不存在，创建钱包对象
 */
func NewWalletsYS(nodeID string) *WalletsYS {
	walletsFile := fmt.Sprintf(walletsFileYS,nodeID)
	//step1：钱包文件不存在
	if _, err := os.Stat(walletsFile); os.IsNotExist(err) {
		fmt.Println("钱包文件不存在。。。")
		wallets := &WalletsYS{}
		wallets.WalletMapYS = make(map[string]*WalletYS)
		return wallets
	}

	wsBytes, err := ioutil.ReadFile(walletsFile)
	if err != nil {
		log.Panic(err)
	}

	gob.Register(elliptic.P256())
	var wallets WalletsYS

	reader := bytes.NewReader(wsBytes)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&wallets)
	if err != nil {

		log.Panic(err)
	}

	return &wallets
}

/*
	创建新wallet并保存
 */
func (ws *WalletsYS) CreateWalletYS(nodeID string)  {
	wallet := NewWalletYS()
	address := wallet.GetAddressYS()

	fmt.Printf("创建的钱包地址：%s\n",address)

	ws.WalletMapYS[string(address)] =wallet

	ws.saveFileYS(nodeID)
}

func (ws *WalletsYS) saveFileYS (nodeID string) {
	//组合文件名
	walletsFile := fmt.Sprintf(walletsFileYS,nodeID)
	//序列化ws对象
	wsBytes := gobEncodeWithRegister(ws, elliptic.P256())
	//将序列化后的ws对象存入文件
	err := ioutil.WriteFile(walletsFile, wsBytes, 0644)
	if err != nil {
		log.Panic(err)
	}
}
