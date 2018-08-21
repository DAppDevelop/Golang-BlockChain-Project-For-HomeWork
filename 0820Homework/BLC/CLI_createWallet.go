package BLC

import "fmt"

func (cli *CLIYS) CreateWalletYS() {
	wallets := NewWalletsYS()
	wallets.CreateWalletYS()
	fmt.Println("钱包：", wallets.WalletMapYS)
}
