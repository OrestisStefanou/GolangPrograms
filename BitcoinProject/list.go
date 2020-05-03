package main

import "fmt"

type listNode struct {
	data *bitcoinNode
	next *listNode
}

type list struct {
	head *listNode
}

//Insert at the head of the list
func (l *list) insert(bNode *bitcoinNode) {
	if l.head == nil { //IF LIST IS EMPTY
		l.head = new(listNode)
		l.head.data = bNode
		l.head.next = nil
		return
	}
	var newnode *listNode = new(listNode)
	newnode.data = bNode
	newnode.next = l.head
	l.head = newnode
}

func (l *list) print() {
	var tempnode *listNode = l.head
	for tempnode != nil {
		fmt.Printf("Bitcoin node owener:%s,value:%d\n", tempnode.data.owner, tempnode.data.value)
		tempnode = tempnode.next
	}

}

//Get the total value of Leaf nodes that are in the list
func (l *list) getTotal() int {
	var tempnode *listNode = l.head
	var sum int = 0
	for tempnode != nil {
		if tempnode.data.left == nil && tempnode.data.right == nil { //if is a leaf node in the Bitcoin Tree
			sum = sum + tempnode.data.value
		}
		tempnode = tempnode.next
	}
	return sum
}

//Get the total value of <k> nodes that are in the list
//if k=-1 return the total earnings
func (l *list) getEarnings(k int, id string) int {
	var tempnode *listNode = l.head
	var sum int = 0
	var counter int = 0
	for tempnode != nil {
		if counter == k {
			break
		}
		sum = sum + tempnode.data.value
		counter = counter + 1
		tempnode = tempnode.next
	}
	return sum
}
