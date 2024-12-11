package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"test/location-forwarder-demo/app/businesses"
	"test/location-forwarder-demo/distribution/interceptors"
	invoker "test/location-forwarder-demo/distribution/invokers"
	"test/shared"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	locationForwarder := interceptors.NewLocationForwarder()

	fmt.Print("Type a port: ")
	scanner.Scan()
	calcPort, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatalf("error reading input 'b': %s", err)
	}

	inv := invoker.NewInvoker(shared.LocalHost, calcPort, &locationForwarder, true)
	go listenCommand(scanner, &inv)

	inv.Invoke()
}

func listenCommand(scanner *bufio.Scanner, invoker *invoker.Invoker) {

	fmt.Println("Type 'exit' to finish server at any moment...")

	for {
		fmt.Println("--------------------------")
		fmt.Print("Type 1 to add or 2 to remove an calculator object: ")
		scanner.Scan()

		switch scanner.Text() {
		case "exit":
			os.Exit(1)
		case "1":
			fmt.Print("Type the calculator name to add: ")
			scanner.Scan()
			invoker.AddLocalObject(scanner.Text(), businesses.Calculadora{})
			fmt.Println("New calculator created.")
		case "2":
			fmt.Print("Type the calculator name to remove: ")
			scanner.Scan()
			calculatorToRemove := scanner.Text()
			fmt.Print("Type the new port: ")
			scanner.Scan()
			newPort, _ := strconv.Atoi(scanner.Text())
			invoker.RemoveLocalObject(calculatorToRemove, shared.IOR{Host: shared.LocalHost, Port: newPort})
			fmt.Println("Calculator migrated from server.")
		default:
			fmt.Println("Unknown command")
		}
	}
}
