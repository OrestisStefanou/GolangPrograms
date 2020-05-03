package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var pipe1 = "/home/orestis/Desktop/pipe1.log"
var pipe2 = "/home/orestis/Desktop/pipe2.log"

//Create the communication pipes
func (network *bitcoinNetwork) startNetwork() {
	err := syscall.Mkfifo(pipe1, 0666) //Server creates the pipe for reading
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

	err = syscall.Mkfifo(pipe2, 0666) //Server creates the pipe for writing
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

}

func (network *bitcoinNetwork) getRequest() {

	var line string
	for {
		file1, err := os.OpenFile(pipe1, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Fatal("Open named pipe file error:", err)
		}
		scanner := bufio.NewScanner(file1)
		for scanner.Scan() {
			line = scanner.Text()
			//		file2.WriteString("Server says hello")
		}
		if err = scanner.Err(); err != nil {
			fmt.Println("error occured")
		}
		file1.Close()
		if line == "end" {
			break
		}
		go network.handleRequest(line)
	}

	os.Remove(pipe1)
	os.Remove(pipe2)
}

func (network *bitcoinNetwork) handleRequest(request string) {
	requestArgs := strings.Split(request, " ")

	if requestArgs[0] == "/requestTransaction" {
		//HANDLE THE TRANSACTION
		if len(requestArgs) != 4 {
			fmt.Println("Wrong usage")
			return
		}
		var senderID, receiverID string
		senderID = requestArgs[1]
		receiverID = requestArgs[2]
		amount, err := strconv.Atoi(requestArgs[3])
		if err != nil {
			fmt.Println("Wrong usage")
			return
		}
		network.requestTransaction(senderID, receiverID, amount, true)
		return
	}

	if requestArgs[0] == "/findEarnings" {
		if len(requestArgs) != 3 {
			fmt.Println("Wrong usage")
			return
		}
		var walletID string = requestArgs[1]
		k, err := strconv.Atoi(requestArgs[2])
		if err != nil {
			fmt.Println("Wrong usage")
			return
		}
		network.findEarnings(walletID, k)
		return
	}

	if requestArgs[0] == "/findPayments" {
		//Handle the request
		if len(requestArgs) != 3 {
			fmt.Println("Wrong usage")
			return
		}
		var walletID string = requestArgs[1]
		k, err := strconv.Atoi(requestArgs[2])
		if err != nil {
			fmt.Println("Wrong usage")
			return
		}
		network.findPayments(walletID, k)
		return
	}

	if requestArgs[0] == "/walletStatus" {
		//Handle the request
		if len(requestArgs) != 2 {
			fmt.Println("Wrong usage")
			return
		}
		var walletID string = requestArgs[1]
		network.walletStatus(walletID)
		return
	}

	if requestArgs[0] == "/bitcoinStatus" {
		//Handle the request
		if len(requestArgs) != 2 {
			fmt.Println("Wrong usage")
			return
		}
		var bitcoinID string = requestArgs[1]
		network.bitcoinStatus(bitcoinID)
		return
	}

	if requestArgs[0] == "/traceCoin" {
		//Handle the request
		if len(requestArgs) != 2 {
			fmt.Println("Wrong usage")
			return
		}
		var bitcoinID string = requestArgs[1]
		network.tracecoin(bitcoinID)
		return
	}

	if requestArgs[0] == "/showWallets" {
		if len(requestArgs) != 1 {
			fmt.Println("Wrong usage")
			return
		}
		network.showWallets()
		return
	}
	//Wrong usage,send appropriate message
	file2, err := os.OpenFile(pipe2, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}
	file2.WriteString("Wrong usage.Type /help for more info\n")
	file2.Close()
}
