package BLC

import "fmt"

func (cli *CLIYS) GetAddressListYS(nodeID string) {
	fmt.Println("打印所有的钱包地址。。")
	wallets := NewWalletsYS(nodeID)
	for address, _ := range wallets.WalletMapYS {
		fmt.Println("address: ", address)
	}
}
