package BLC

import "fmt"

func (cli *CLIYS) CreateWalletYS(nodeID string) {
	wallets := NewWalletsYS(nodeID)
	wallets.CreateWalletYS(nodeID)
	fmt.Println("钱包：", wallets.WalletMapYS)
}
