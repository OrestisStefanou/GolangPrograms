package main

import "fmt"

type tnode struct {
	data  int
	left  *tnode
	right *tnode
}

//Add a node to the tree
func addtree(number int, p *tnode) *tnode {
	if p == nil {
		p = new(tnode)
		p.data = number
		p.left = nil
		p.right = nil

	}
	if number < p.data {
		p.left = addtree(number, p.left)

	}
	if number > p.data {
		p.right = addtree(number, p.right)

	}
	return p //TRY IF IT WORKS WITHOUT THIS
}

//Print the nodes of the tree
func nodesprint(p *tnode) {
	if p != nil {
		nodesprint(p.left)
		fmt.Println(p.data)
		nodesprint(p.right)
	}
}

//Return the depth of the tree
func treedepth(p *tnode) int {
	var n1, n2 int
	if p == nil {
		return 0
	}
	n1 = treedepth(p.left)
	n2 = treedepth(p.right)
	if n1 > n2 {
		return n1 + 1
	}
	return n2 + 1
}

//Check if data exists in the tree
func treesearch(p *tnode, number int) int {
	if p == nil { //number not in the tree
		return 0
	}
	if p.data == number { //found the number
		return 1
	}
	if number > p.data { //if number is bigger than current node go to right subtree
		return treesearch(p.right, number)
	}
	//Else go to left subtree
	return treesearch(p.left, number)
}

func main() {
	var root *tnode
	root = addtree(10, root)
	root = addtree(20, root)
	root = addtree(15, root)
	root = addtree(7, root)
	root = addtree(3, root)
	nodesprint(root)
	root = addtree(30, root)
	root = addtree(40, root)
	fmt.Println("After inserting 30")
	nodesprint(root)
	fmt.Println("Tree depth is:")
	fmt.Println(treedepth(root))
	var found int
	found = treesearch(root, 10)
	if found == 1 {
		fmt.Println("Number 10 found")
	} else {
		fmt.Println("Number 10 not in the tree")
	}
	found = treesearch(root, 100)
	if found == 1 {
		fmt.Println("Number 100 found")
	} else {
		fmt.Println("Number 100 not in the tree")
	}

}
