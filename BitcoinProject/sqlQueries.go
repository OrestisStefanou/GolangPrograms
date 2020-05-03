package main

import (
	"database/sql"
	"fmt"
	"log"
)

func dbTransactionInsert(bitcoinID string, senderID string, receiverID string, value int) {
	db, err := sql.Open("mysql", "root:19971948Oo@/Bitcoin") //Connect to the database
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Transactions(BitcoinID,senderID,receiverID,value)VALUES(?,?,?,?)")
	if err != nil {
		fmt.Println("Error during database transaction..Aborting")
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(bitcoinID, senderID, receiverID, value)
	if err != nil {
		fmt.Println("Error during database transaction..Aborting")
		return
	}
}
