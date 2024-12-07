package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	calculadorainvoker "test/myrpc/distribution/invokers/calculadora"
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
	calcInv := calculadorainvoker.New(shared.LocalHost, calcPort)

	// Register services in Naming
	naming.Bind("Calculadora", shared.NewIOR(calcInv.Ior.Host, calcInv.Ior.Port))

	// Invoke services
	calcInv.Invoke()
}
