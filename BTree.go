package main

import (
	"fmt"
)

const (
	max int = 4 //Maximum number of keys in node
	min int = 2 //Minimum number of keys in node
)

type keyType int

type treeEntry struct {
	key  keyType //The key that the tree will be sorted
	data interface{}
}

type treeNode struct {
	count  int //denotes how many keys there are in the node
	entry  [max + 1]treeEntry
	branch [max + 1]*treeNode
}

/* createTree: create a B-tree.
   Pre: None.
   Post: An empty B-tree has been created to which root points. */
func createTree(root **treeNode) {
	*root = nil
}

/* Preorder: print each node of the B-tree in preorder.
   Pre: The B-tree to which root points has been created.
   Post: Each node of the B-tree has been printed in preorder. */
func preOrder(root *treeNode) {
	if root != nil {
		for i := 1; i <= root.count; i++ {
			fmt.Printf("%d\n", root.entry[i].key)
		}
		fmt.Println("End of node")

		for i := 0; i < root.count+1; i++ {
			preOrder(root.branch[i])
		}
	}
}

/* SearchNode: searches keys in node for target.
   Pre: target is a key and current points to a node of a B-tree.
   Post: Searches keys in node for target; returns location pos of target, or
   branch on which to continue search.*/
func searchNode(target keyType, current *treeNode, pos *int) bool {
	if target < current.entry[1].key { //Take the leftmost branch
		*pos = 0
		return false
	}
	for *pos = current.count; target < current.entry[*pos].key && *pos > 1; (*pos)-- {

	}
	return (target == current.entry[*pos].key)
}

/* SearchTree: traverse the B-tree looking for target.
   Pre: If the key target is present in the B-tree, then the return value points to
   the node containing target in position targetpos. Otherwise, the return value is
   NULL and targetpos is undefined.
   Uses: SearchTree recursively, SearchNode. */
func searchTree(target keyType, root *treeNode, targetpos *int) *treeNode {
	if root == nil {
		return nil
	}
	if searchNode(target, root, targetpos) {
		return root
	}
	return searchTree(target, root.branch[*targetpos], targetpos)
}

/* PushIn: inserts a key into a node.
   Pre: medentry belongs at index pos in node *current, which has room for it.
   Post: Inserts key medentry and pointer medright into *current at index pos. */
func pushIn(medentry treeEntry, medright *treeNode, current *treeNode, pos int) {
	var i int //Index to move keys to make room for medentry

	for i = current.count; i > pos; i-- {
		//Shift all keys and branches to the right
		current.entry[i+1] = current.entry[i]
		current.branch[i+1] = current.branch[i]
	}
	current.entry[pos+1] = medentry
	current.branch[pos+1] = medright
	current.count++
}

/* Split: splits a full node.
   Pre: medentry belongs at index pos of node *current which is full.
   Post: Splits node *current with entry medentry and pointer medright at index pos
         into nodes *current and *newright with median entry newmedian.
   Uses: PushIn */
func split(medentry treeEntry, medright *treeNode, current *treeNode, pos int, newmedian *treeEntry, newright **treeNode) {
	var i int      //used for copying from *current to new node
	var median int //median position in the combined,overfull node

	if pos <= min { //Find splitting point.Determine if key goes to left or right half
		median = min
	} else {
		median = min + 1
	}

	//Get a new node and put it on the right
	*newright = new(treeNode)
	for i = median + 1; i <= max; i++ { //Move half the keys to the right node
		(*newright).entry[i-median] = current.entry[i]
		(*newright).branch[i-median] = current.branch[i]
	}
	(*newright).count = max - median
	current.count = median

	if pos <= min { //Push in the new key
		pushIn(medentry, medright, current, pos)
	} else {
		pushIn(medentry, medright, *newright, pos-median)
	}
	*newmedian = current.entry[current.count]
	(*newright).branch[0] = current.branch[current.count]
	current.count--
}

/* PushDown: recursively move down tree searching for newentry.
   Pre: newentry belongs in the subtree to which current points.
   Post: newentry has been inserted into the subtree to which current points; if TRUE
   is returned, then the height of the subtree has grown, and medentry needs
   to be reinserted higher in the tree, with subtree medright on its right.
   Uses: PushDown recursively, SearchNode, Split, PushIn. */
func pushDown(newentry treeEntry, current *treeNode, medentry *treeEntry, medright **treeNode) bool {
	var pos int //branch on wich to continues the search

	if current == nil {
		*medentry = newentry
		*medright = nil
		return true
	}
	//Search the current node
	if searchNode(newentry.key, current, &pos) {
		fmt.Println("Inserting duplicate tree in B-tree")
		return false
	}
	if pushDown(newentry, current.branch[pos], medentry, medright) {
		if current.count < max { //Reinsert median key\
			pushIn(*medentry, *medright, current, pos)
			return false
		}
		split(*medentry, *medright, current, pos, medentry, medright)
		return true
	}
	return false
}

/* InsertTree: Inserts entry into the B-tree.
   Pre: The B-tree to which root points has been created, and no entry in the B-tree
   has key equal to newentry key.
   Post: newentry has been inserted into the B-tree, the root is returned.
   Uses: PushDown */
func insertTree(newentry treeEntry, root *treeNode) *treeNode {
	var medentry treeEntry //node to be reinserted as new root
	var medright *treeNode //subtree on right of medentry
	var newroot *treeNode  //used when the heigth of the tree increases

	if pushDown(newentry, root, &medentry, &medright) {
		//Tree grows in height.Make a new root
		newroot = new(treeNode)
		newroot.count = 1
		newroot.entry[1] = medentry
		newroot.branch[0] = root
		newroot.branch[1] = medright
		return newroot
	}
	return root
}

/* Combine: combine adjacent nodes.
   Pre: current points to a node in a B-tree with entries in the branches pos and
   pos-1, with too few to move entries.
   Post: The nodes at branches pos-1 and pos have been combined into one node,
   which also includes the entry formerly in *current at index pos. */
func combine(current *treeNode, pos int) {
	var c int
	var right *treeNode
	var left *treeNode
	right = current.branch[pos]
	left = current.branch[pos-1] //work with the left node
	left.count++                 //insert the key from the parent
	left.entry[left.count] = current.entry[pos]
	left.branch[left.count] = right.branch[0]
	for c = 1; c <= right.count; c++ { //Insert all jeys from right node
		left.count++
		left.entry[left.count] = right.entry[c]
		left.branch[left.count] = right.branch[c]
	}
	for c = pos; c < current.count; c++ { //Delete key from parent node
		current.entry[c] = current.entry[c+1]
		current.branch[c] = current.branch[c+1]
	}
	current.count--
}

/* MoveLeft: move a key to the left.
Pre: current points to a node in a B-tree with entries in the branches pos and
pos-1, with too few in branch pos-1.
Post: The leftmost entry from branch pos has moved into *current, which has sent
an entry into the branch pos-1 */
func moveLeft(current *treeNode, pos int) {
	var c int
	var t *treeNode
	t = current.branch[pos-1] //Move key from parent into left node
	t.count++
	t.entry[t.count] = current.entry[pos]
	t.branch[t.count] = current.branch[pos].branch[0]
	t = current.branch[pos] //Move key from right into parent
	current.entry[pos] = t.entry[1]
	t.branch[0] = t.branch[1]
	t.count--
	for c = 1; c <= t.count; c++ {
		//shift all keys in the right node one posiotion leftward
		t.entry[c] = t.entry[c+1]
		t.branch[c] = t.branch[c+1]
	}
}

/* MoveRight: move a key to the right.
   Pre: current points to a node in a B-tree with entries in the branches pos and
   pos-1, with too few entries in branch pos.
   Post: The rightmost entry from branch pos-1 has moved into *current, which has
   sent an entry into the branch pos */
func moveRight(current *treeNode, pos int) {
	var c int
	var t *treeNode
	t = current.branch[pos]
	for c = t.count; c > 0; c-- {
		//shift all keys in the right node one position
		t.entry[c+1] = t.entry[c]
		t.branch[c+1] = t.branch[c]
	}
	t.branch[1] = t.branch[0] //move key from parent to right node
	t.count++
	t.entry[1] = current.entry[pos]
	t = current.branch[pos-1] //move last key of left node into parent
	current.entry[pos] = t.entry[t.count]
	current.branch[pos].branch[0] = t.branch[t.count]
	t.count--
}

func main() {
	var r *treeNode
	var e treeEntry
	createTree(&r)
	for i := 0; i < 1000; i++ {
		e.key = keyType(i)
		e.data = i
		r = insertTree(e, r)
	}
	fmt.Println("Nodes of the tree in preorder")
	preOrder(r)

	//Test the search for a target
	var pos int
	s := searchTree(20, r, &pos)
	if s != nil {
		fmt.Printf("The key is in position %d of the node with the following elements:\n", pos)
		for i := 1; i <= s.count; i++ {
			fmt.Printf("%d\n", s.entry[i].key)
		}
	} else {
		fmt.Println("Key not found")
	}

}
