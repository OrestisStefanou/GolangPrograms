package main

import "fmt"

type listnode struct {
	data int
	next *listnode
}

type list struct {
	head *listnode
	tail *listnode
}

func initlist(listptr *list) {
	listptr.head = nil
	listptr.tail = nil
}

//Insert at the end of the list
func append(listptr *list, data int) *list {
	if listptr.head == nil { //List is empty
		listptr.head = new(listnode)
		listptr.head.data = data
		listptr.head.next = nil
		listptr.tail = listptr.head
		return listptr
	}
	var newnode *listnode = new(listnode)
	newnode.data = data
	newnode.next = nil
	listptr.tail.next = newnode
	listptr.tail = newnode
	return listptr

}

//Insert at the start of the list
func listpush(listptr *list, data int) *list {
	if listptr.head == nil { //List is empty
		listptr.head = new(listnode)
		listptr.head.data = data
		listptr.head.next = nil
		listptr.tail = listptr.head
		return listptr
	}
	var newnode *listnode = new(listnode)
	newnode.data = data
	newnode.next = listptr.head
	listptr.head = newnode
	return listptr
}

func printlist(listptr *list) {
	var tempptr *listnode = listptr.head
	for tempptr != nil {
		fmt.Print(tempptr.data)
		fmt.Print(" ")
		tempptr = tempptr.next
	}
}

func main() {
	var mylist *list = new(list)
	initlist(mylist)
	mylist = append(mylist, 1)
	mylist = append(mylist, 2)
	mylist = append(mylist, 3)
	mylist = append(mylist, 4)
	mylist = append(mylist, 5)
	mylist = append(mylist, 6)
	mylist = listpush(mylist, 10)
	printlist(mylist)
}
