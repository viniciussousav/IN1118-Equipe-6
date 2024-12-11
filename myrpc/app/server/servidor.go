package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"test/myrpc/distribution/interceptors"
	invoker "test/myrpc/distribution/invokers"
	"test/shared"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	// Create instance of location forwarder
	locationForwarder := interceptors.NewLocationForwarder()

	// Create instance of invokers
	calcPort := getPort(scanner)
	go listenExitCommand(scanner)

	inv := invoker.NewInvoker(shared.LocalHost, calcPort, &locationForwarder)

	// Invoke services
	inv.Invoke()
}

func getPort(scanner *bufio.Scanner) int {
	fmt.Print("Type a port: ")
	scanner.Scan()
	calcPort, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatalf("error reading input 'b': %s", err)
	}
	return calcPort
}

func listenExitCommand(scanner *bufio.Scanner) {
	for {
		fmt.Println("Type 'exit' to finish server at any moment...")
		scanner.Scan()
		if scanner.Text() == "exit" {
			os.Exit(1)
		}
	}
}
