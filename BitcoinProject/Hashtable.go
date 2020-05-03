package main

type hashtable struct {
	ht map[string]*list
}

func (table *hashtable) initHashtable() {
	table.ht = make(map[string]*list)
}

//CREATE NEW KEY:VALUE ENTRY IN THE HASHTABLE
func (table *hashtable) setKey(key string) {
	var tableList *list = new(list)
	tableList.head = nil
	table.ht[key] = tableList
}

//ADD A NODE TO THE LIST WITH KEY key
func (table *hashtable) insert(key string, bNode *bitcoinNode) {
	value, hasKey := table.ht[key]
	if hasKey { //If the key exists in the hashtable
		value.insert(bNode) //Insert the bitcoin node in the list
	} else { //Key does not exist in the hashtable
		table.setKey(key) //we create a new entry
		value = table.ht[key]
		value.insert(bNode)
	}
}

//Returns the total value of a Wallet
func (table *hashtable) getWalletValue(key string) int {
	value, hasKey := table.ht[key]
	if hasKey { //Key exists in the hashtable
		var totalvalue int = value.getTotal()
		return totalvalue
	}
	return -1
}

//TO USE IN FINDEARNINGS
func (table *hashtable) getEarnings(key string, k int) int {
	value, hasKey := table.ht[key]
	if hasKey {
		var earnings int = value.getEarnings(k, key)
		return earnings
	}
	return -1
}
