package main

import (
	"bufio"
	"fmt"
	"os"
)

func printHelp() {
	fmt.Println("Request examples:")
	fmt.Println("/requestTransaction <senderWalletID> <receiverWalletID> <amount> ")
	fmt.Println("/findEarnings <walletID> <number of transactions(-1 to find total earnings)>")
	fmt.Println("/findPayents <walletID> <number of transactions(-1 to find total amount of payments)>")
	fmt.Println("/walletStatus <walletID>")
	fmt.Println("/bitcoinStatus <bitcoinID>")
	fmt.Println("/traceCoin <bitcoinID>")
	fmt.Println("/showWallets")
}

var pipe1File = "/home/orestis/Desktop/pipe1.log"
var pipe2File = "/home/orestis/Desktop/pipe2.log"

func main() {
	var line string
	for {
		f1, err := os.OpenFile(pipe1File, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			fmt.Println("Error occured")
		}

		scanner := bufio.NewScanner(os.Stdin)

		scanner.Scan()
		line = scanner.Text()
		line = line + "\n"
		if line == "/help\n" {
			printHelp()
			f1.Close()
			continue
		}
		f1.WriteString(line)

		f1.Close()
		if line == "end\n" {
			break
		}

		if line == "/help\n" {

		}

		f2, err := os.OpenFile(pipe2File, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			fmt.Println("Error occured")
		}
		scanner2 := bufio.NewScanner(f2)

		for scanner2.Scan() {
			readline := scanner2.Text()
			fmt.Println(readline)
		}
		f2.Close()
	}
}
