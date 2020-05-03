package main

import (
	"sync"
)

var mu sync.Mutex

type queueNode struct {
	request string
	next    *queueNode
}

type queue struct {
	head *queueNode
	tail *queueNode
}

func (q *queue) init() {
	q.head = nil
	q.tail = nil
}

//More than one go routines may try to insert at the same time
func (q *queue) insert(r string) {
	mu.Lock()
	var newnode *queueNode = new(queueNode)
	newnode.request = r
	newnode.next = nil
	if q.head == nil { //Queue is empty
		q.head = newnode
		q.tail = newnode
		mu.Unlock()
		return
	}
	q.tail.next = newnode
	q.tail = newnode
	mu.Unlock()
}

//More than one go routines may try to extract at the same time
func (q *queue) pop() string {
	mu.Lock()
	var request string = q.head.request
	q.head = q.head.next
	mu.Unlock()
	return request
}
