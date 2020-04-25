package main

import "fmt"

type rbNode struct {
	color string //RED OR BLACK
	data  int
	link  [2]*rbNode
}

func createNode(number int) *rbNode {
	var newnode *rbNode
	newnode = new(rbNode)
	newnode.data = number
	newnode.color = "RED"
	newnode.link[0] = nil
	newnode.link[1] = nil
	return newnode
}

func rbInsertion(number int, root *rbNode) *rbNode {
	var stack = [100]*rbNode{} //Stack to keep parent nodes
	var ptr *rbNode
	var newnode *rbNode
	//Pointers used for rotations
	var rotationPtr1 *rbNode
	var rotationPtr2 *rbNode
	///////////////////////
	var linkIndex = [100]int{} //Keep left of right pointer index(0->left,1->right)
	//Counters for stack and link_index arrays
	var counter int = 0
	var index int
	////////////////
	ptr = root
	if root == nil { //If tree is empty
		root = createNode(number)
		return root
	}
	stack[counter] = root
	linkIndex[counter] = 0
	counter = counter + 1
	//Find place to insert the new node
	for ptr != nil {
		if number > root.data {
			index = 1
		} else {
			index = 0
		}
		stack[counter] = ptr
		ptr = ptr.link[index]
		linkIndex[counter] = index
		counter = counter + 1
	}
	//Insert new node
	newnode = createNode(number)
	stack[counter-1].link[index] = newnode
	for counter >= 3 && stack[counter-1].color == "RED" {
		if linkIndex[counter-2] == 0 {
			rotationPtr2 = stack[counter-2].link[1]
			if rotationPtr2 != nil && rotationPtr2.color == "RED" {
				stack[counter-2].color = "RED"
				stack[counter-1].color = "BLACK"
				rotationPtr2.color = "BLACK"
				counter = counter - 2
			} else {
				if linkIndex[counter-1] == 0 {
					rotationPtr2 = stack[counter-1]
				} else {
					rotationPtr1 = stack[counter-1]
					rotationPtr2 = rotationPtr1.link[1]
					rotationPtr1.link[1] = rotationPtr2.link[0]
					rotationPtr2.link[0] = rotationPtr1
					stack[counter-2].link[0] = rotationPtr2
				}
				rotationPtr1 = stack[counter-2]
				rotationPtr1.color = "RED"
				rotationPtr2.color = "BLACK"
				rotationPtr1.link[0] = rotationPtr2.link[1]
				rotationPtr2.link[1] = rotationPtr1
				if rotationPtr1 == root {
					root = rotationPtr2
				} else {
					stack[counter-3].link[linkIndex[counter-3]] = rotationPtr2
				}
				break
			}
		} else {
			rotationPtr2 = stack[counter-2].link[0]
			if rotationPtr2 != nil && rotationPtr2.color == "RED" {
				stack[counter-2].color = "RED"
				stack[counter-1].color = "BLACK"
				rotationPtr2.color = "BLACK"
				counter = counter - 2
			} else {
				if linkIndex[counter-1] == 1 {
					rotationPtr2 = stack[counter-1]
				} else {
					rotationPtr1 = stack[counter-1]
					rotationPtr2 = rotationPtr1.link[0]
					rotationPtr1.link[0] = rotationPtr2.link[1]
					rotationPtr2.link[1] = rotationPtr1
					stack[counter-2].link[1] = rotationPtr2
				}
				rotationPtr1 = stack[counter-2]
				rotationPtr2.color = "BLACK"
				rotationPtr1.color = "RED"
				rotationPtr1.link[1] = rotationPtr2.link[0]
				rotationPtr2.link[0] = rotationPtr1
				if rotationPtr1 == root {
					root = rotationPtr2
				} else {
					stack[counter-3].link[linkIndex[counter-3]] = rotationPtr2
				}
				break
			}
		}
	}
	root.color = "BLACK"
	return root
}

//Check if data exists in the tree
func rbSearch(root *rbNode, number int) int {
	if root == nil { //number not in the tree
		return 0
	}
	if root.data == number { //found the number
		return 1
	}
	if number > root.data { //if number is bigger than current node go to right subtree
		return rbSearch(root.link[1], number)
	}
	//Else go to left subtree
	return rbSearch(root.link[0], number)
}

func rbPrint(root *rbNode) {
	if root != nil {
		rbPrint(root.link[0])
		fmt.Println(root.data)
		rbPrint(root.link[1])
	}
}

func main() {
	var root *rbNode = nil
	root = rbInsertion(5, root)
	root = rbInsertion(10, root)
	root = rbInsertion(2, root)
	root = rbInsertion(3, root)
	root = rbInsertion(1, root)
	root = rbInsertion(8, root)
	root = rbInsertion(20, root)
	rbPrint(root)
	var found int
	found = rbSearch(root, 20)
	if found == 1 {
		fmt.Println("Number 20 is in the tree")
	} else {
		println("Number 20 is not in the tree")
	}
	found = rbSearch(root, 100)
	if found == 1 {
		fmt.Println("Number 100 is in the tree")
	} else {
		fmt.Println("Number 100 is not in the tree")
	}
}
