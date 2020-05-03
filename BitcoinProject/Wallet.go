package main

type wallet struct {
	walletID   string
	publicKey  string
	privateKey string
	totalValue int
}

func createWallet(id string, plkey string, pvkey string) *wallet {
	newWallet := new(wallet)
	newWallet.walletID = id
	newWallet.publicKey = plkey
	newWallet.privateKey = pvkey
	return newWallet
}
