package BLC

// createWallets 创建钱包集合
func (cli *CLI) createWallets() {
	wallets := NewWallets()  // 创建一个集合对象
	wallets.CreateWallet()
}