package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"test/myrpc/distribution/interceptors"
	invoker "test/myrpc/distribution/invokers/calculadora"
	namingproxy "test/myrpc/services/naming/proxy"
	"test/shared"
)

func main() {
	// Obtain proxies
	naming := namingproxy.New(shared.LocalHost, shared.NamingPort)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	calcPort, err := strconv.Atoi(scanner.Text())
	if err != nil {
		log.Fatal(err)
	}

	// Create instance of location forwarder
	locationForwarder := interceptors.NewLocationForwarder()

	// Create instance of invokers
	inv := invoker.NewInvoker(shared.LocalHost, calcPort, &locationForwarder)

	// Register services in Naming
	naming.Bind("Calculadora", shared.NewIOR(shared.LocalHost, calcPort))

	// Invoke services
	inv.Invoke()
}
