package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
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

	// Create instance of invokers
	inv := invoker.NewInvoker(shared.LocalHost, calcPort)

	// Register services in Naming
	naming.Bind("Calculadora", shared.NewIOR(inv.Ior.Host, inv.Ior.Port))

	// Invoke services
	inv.Invoke()
}
