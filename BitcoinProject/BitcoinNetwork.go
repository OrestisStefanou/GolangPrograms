package main

import (
	"database/sql"
	"fmt"
	"os"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

type bitcoinNetwork struct {
	senderHashtable   hashtable               //Sending bitcoin nodes of each user
	receiverHashtable hashtable               //Received bitcoin nodes of each user
	bitcoins          map[string]*bitcoinNode //A hashtable with all the bitcoin 's root
	wallets           map[string]*wallet
}

//REQUEST A TRANSACTION FROM THE NETWORK
//If flag=true we update the database
func (network *bitcoinNetwork) requestTransaction(senderID string, receiverID string, amount int, flag bool) {
	var file2 *os.File
	//OPEN PIPE
	if flag {
		file2, _ = os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
		if file2 == nil {
			fmt.Println("Error occured during opening pipe")
			return
		}
	}
	/////
	fmt.Printf("WalletID:%s requesting to send %d to WalletID:%s\n", senderID, amount, senderID)
	//CHECK IF SENDER HAS ENOUGH BITCOINS
	var senderWalletValue int = network.wallets[senderID].totalValue
	if senderWalletValue < amount {
		fmt.Printf("WalletID:%s value is not enough\n", senderID)
		if flag {
			message := fmt.Sprintf("WalletID:%s value is not enough\n", senderID)
			file2.WriteString(message)
			file2.Close()
		}
		return
	}
	var newSenderNode, newReceiverNode *bitcoinNode                       //To get the 2 new nodes that will be created from the transaction
	var senderBitcoinNodes *list = network.receiverHashtable.ht[senderID] //All the bitcoin nodes that the sender received
	var listHead *listNode = senderBitcoinNodes.head                      //The first node from senderBitcoinNodes list
	var amountLeft int = amount
	for amountLeft > 0 { //While the amount left to send is bigger than zero
		for listHead.data.isLeaf() != true { //Find a node that is a leaf so the sender can spend it
			listHead = listHead.next
		}
		if listHead.data.value >= amountLeft { //if the value of this node is >= to the amount left to send
			newSenderNode, newReceiverNode = makeTransaction(listHead.data, amountLeft, senderID, receiverID) //make the transaction
			if flag {
				dbTransactionInsert(listHead.data.bitcoinID, senderID, receiverID, amountLeft) //Update the databse
			}
			//UPDATE THE HASHTABLES
			if newSenderNode.value == 0 {
				//We will not insert newSenderNode in the ReceiversHashtable because it has no value
				network.receiverHashtable.insert(receiverID, newReceiverNode)
				network.senderHashtable.insert(senderID, newReceiverNode)
				fmt.Println("Requested transaction completed successfully")
				break
			} else {
				network.receiverHashtable.insert(receiverID, newReceiverNode)
				network.receiverHashtable.insert(senderID, newSenderNode)
				network.senderHashtable.insert(senderID, newReceiverNode)
				fmt.Println("Requested transaction completed successfully")
				break
			}
		} else { //The value of the node is not enough to complete the request
			newSenderNode, newReceiverNode = makeTransaction(listHead.data, listHead.data.value, senderID, receiverID) //make the transaction
			if flag {
				dbTransactionInsert(listHead.data.bitcoinID, senderID, receiverID, listHead.data.value) //Update the databse
			}
			//We will not insert newSenderNode in the ReceiversHashtable because it has no value
			network.receiverHashtable.insert(receiverID, newReceiverNode)
			network.senderHashtable.insert(senderID, newReceiverNode)
			amountLeft = amountLeft - listHead.data.value //Reduce amount left
			listHead = listHead.next                      //Get the next Bitcoin Node
		}
	}
	//Update the wallet values
	network.wallets[senderID].totalValue = network.receiverHashtable.getWalletValue(senderID)
	network.wallets[receiverID].totalValue = network.receiverHashtable.getWalletValue(receiverID)
	if flag {
		file2.WriteString("Requested transaction completed successfully")
		file2.Close()
	}
}

func (network *bitcoinNetwork) initNetwork() {
	//Initialize data structures
	network.senderHashtable.initHashtable()
	network.receiverHashtable.initHashtable()
	network.bitcoins = make(map[string]*bitcoinNode)
	network.wallets = make(map[string]*wallet)
	//Fill data structures
	db, err := sql.Open("mysql", "root:19971948Oo@/Bitcoin") //Connect to the database
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var bnode *bitcoinNode //Pointer to a structure of a bitcoin node
	//Get bitcoins data
	rows, err := db.Query("SELECT * FROM Bitcoins")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		bnode = new(bitcoinNode)
		if err := rows.Scan(&bnode.bitcoinID, &bnode.owner, &bnode.value); err != nil {
			log.Fatal(err)
		}
		network.bitcoins[bnode.bitcoinID] = bnode            //Add the bitcoin to the hashtable
		network.receiverHashtable.insert(bnode.owner, bnode) //Add to the owener received bitcoins the bitcoin node
	}
	//Get wallets data
	rows, err = db.Query("SELECT * FROM Wallets")
	if err != nil {
		log.Fatal(err)
	}
	var newWallet *wallet
	for rows.Next() {
		newWallet = new(wallet)
		if err := rows.Scan(&newWallet.walletID, &newWallet.publicKey, &newWallet.publicKey); err != nil {
			log.Fatal(err)
		}
		network.wallets[newWallet.walletID] = newWallet //Add the wallet to wallets hashtable
	}
	//SET Wallets total value
	var totalValue int
	for key, value := range network.wallets {
		totalValue = network.receiverHashtable.getWalletValue(key) //Get the total value of the wallet
		value.totalValue = totalValue
	}

	//Make the transactions
	var tInfo transactionInfo
	rows, err = db.Query("SELECT senderID,receiverID,value FROM Transactions")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		if err := rows.Scan(&tInfo.senderID, &tInfo.receicerID, &tInfo.value); err != nil {
			log.Fatal(err)
		}
		network.requestTransaction(tInfo.senderID, tInfo.receicerID, tInfo.value, false)
	}
	//Update wallets
	for key, value := range network.wallets {
		totalValue = network.receiverHashtable.getWalletValue(key) //Get the total value of the wallet
		value.totalValue = totalValue
	}
}

//Get the total earnings of last <number> transaction of walletID
//If number=-1 return the total earnings from all the transactions
//Returns -1 in case of error
func (network *bitcoinNetwork) findEarnings(walletID string, number int) {
	var earnings int = network.receiverHashtable.getEarnings(walletID, number)
	var message string
	file2, err := os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}
	if earnings == -1 { //Error occured
		file2.WriteString("Requested WalletID has no earnings\n")
		file2.Close()
		return
	}
	if number == -1 {
		message = fmt.Sprintf("Total earnings of WalletID:%s are:%d\n", walletID, earnings)
	} else {
		message = fmt.Sprintf("Total earnings from last %d transactions of WalletID:%s are:%d\n", number, walletID, earnings)
	}
	file2.WriteString(message)
	file2.Close()
}

//Get the total amount of payments of last <number> transaction of walletID
//If number=-1 return the total payments from all the transactions
//Returns -1 in case of error
func (network *bitcoinNetwork) findPayments(walletID string, number int) {
	var payments int = network.senderHashtable.getEarnings(walletID, number)
	var message string
	file2, err := os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
		return
	}
	if payments == -1 { //Error occured
		file2.WriteString("Requested WalletID has no earnings\n")
		file2.Close()
		return
	}
	if number == -1 {
		message = fmt.Sprintf("Total payments of WalletID:%s are:%d\n", walletID, payments)
	} else {
		message = fmt.Sprintf("Total payments from last %d transactions of WalletID:%s are:%d\n", number, walletID, payments)
	}
	file2.WriteString(message)
	file2.Close()
}

//Get the value of wallet walletID
func (network *bitcoinNetwork) walletStatus(walletID string) {
	var totalvalue int = network.wallets[walletID].totalValue
	file2, err := os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
		return
	}
	message := fmt.Sprintf("WalletID:%s total value is:%d\n", walletID, totalvalue)
	file2.WriteString(message)
	file2.Close()
}

//Show the history of transactions of bitcoin bitcoinID
func (network *bitcoinNetwork) tracecoin(bitcoinID string) {
	printBitcoinNodes(network.bitcoins[bitcoinID])
}

//Show bitcoin's current Owners
func (network *bitcoinNetwork) bitcoinStatus(bitcoinID string) {
	printBitcoinOwners(network.bitcoins[bitcoinID])
}

//Show total value of every wallet in the Network
func (network *bitcoinNetwork) showWallets() {
	file2, err := os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
		return
	}
	fmt.Println("Total value of each wallet:")
	file2.WriteString("Total value of each wallet:\n")
	for _, value := range network.wallets {
		fmt.Printf("WalletID:%s,value:%d\n", value.walletID, value.totalValue)
		message := fmt.Sprintf("WalletID:%s,Value:%d\n", value.walletID, value.totalValue)
		file2.WriteString(message)
	}
	file2.Close()
}

func main() {
	var network bitcoinNetwork
	network.startNetwork()
	network.initNetwork()
	network.getRequest()

}
