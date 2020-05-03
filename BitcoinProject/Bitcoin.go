package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

//Structures to save the data we fetch from the database and to help us initialize the network
type bitcoinInfo struct {
	bitcoinID  string
	firstOwner string
	value      int
}

type transactionInfo struct {
	transactionID string
	bitcoinID     string
	senderID      string
	receicerID    string
	value         int
}

//MALLON EN ASKOPA
///////////////////////////////////////////////////

type bitcoinNode struct {
	bitcoinID string
	owner     string
	value     int
	timeDate  time.Time
	left      *bitcoinNode
	right     *bitcoinNode
}

//Create a node of the bitcoin tree
func createBitcoinNode(owner string, value int, id string) *bitcoinNode {
	var newnode *bitcoinNode = new(bitcoinNode)
	newnode.bitcoinID = id
	newnode.owner = owner
	newnode.value = value
	newnode.timeDate = time.Now()
	newnode.left = nil
	newnode.right = nil
	return newnode
}

//Returns true if the node in the bitcoin tree is a leaf
func (bnode *bitcoinNode) isLeaf() bool {
	if bnode.right == nil && bnode.left == nil {
		return true
	}
	return false
}

//Create the two new nodes of the bitcoin tree and return the two new pointers(sender argument maybe is useless)
func makeTransaction(bNode *bitcoinNode, value int, sender string, receiver string) (*bitcoinNode, *bitcoinNode) {
	//Check if the node is a leaf
	if bNode.left != nil || bNode.right != nil {
		fmt.Println("Bitcoin node is not a leaf in the tree")
		return nil, nil
	}

	//Check if owner is the same with the sender
	if bNode.owner != sender {
		fmt.Println("Error during transasction.Bitcoin's owner is not the sender")
		return nil, nil
	}
	//Check if value of the node is bigger than the requesting transaction value
	if bNode.value < value {
		fmt.Println("Bitcoin's value is smaller than requested value")
		return nil, nil
	}
	//Make the transaction
	var newSenderNode *bitcoinNode = createBitcoinNode(bNode.owner, bNode.value-value, bNode.bitcoinID) //right child
	var newReceiverNode *bitcoinNode = createBitcoinNode(receiver, value, bNode.bitcoinID)              //left child

	bNode.right = newSenderNode
	bNode.left = newReceiverNode

	return newSenderNode, newReceiverNode
}

func printBitcoinNodes(bNode *bitcoinNode) {
	if bNode != nil && bNode.left != nil && bNode.right != nil {
		file2, err := os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatal("Open named pipe file error:", err)
			return
		}
		fmt.Printf("Wallet with id:%s transfered %d value of bitcoin to %s\n", bNode.owner, bNode.left.value, bNode.left.owner)
		message := fmt.Sprintf("Wallet with id:%s transfered %d value of bitcoin to %s\n", bNode.owner, bNode.left.value, bNode.left.owner)
		file2.WriteString(message)
		file2.Close()
		printBitcoinNodes(bNode.left)
		printBitcoinNodes(bNode.right)
	}
}

func printBitcoinOwners(bNode *bitcoinNode) {
	if bNode == nil {
		return
	}

	if bNode.left == nil && bNode.right == nil && bNode.value > 0 {
		file2, err := os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatal("Open named pipe file error:", err)
			return
		}
		fmt.Printf("WalletID:%s,value:%d\n", bNode.owner, bNode.value)
		message := fmt.Sprintf("WalletID:%s,value:%d\n", bNode.owner, bNode.value)
		file2.WriteString(message)
		file2.Close()
	}

	if bNode.left != nil {
		printBitcoinOwners(bNode.left)
	}

	if bNode.right != nil {
		printBitcoinOwners(bNode.right)
	}
}
