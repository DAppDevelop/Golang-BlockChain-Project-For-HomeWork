package BLC

import "fmt"

func (cli *CLIYS) GetAddressListYS() {
	fmt.Println("打印所有的钱包地址。。")
	wallets := NewWalletsYS()
	for address := range wallets.WalletMapYS {
		fmt.Println("address: ", address)
	}
}
