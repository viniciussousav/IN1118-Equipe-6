package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"test/myrpc/distribution/proxies"
	namingproxy "test/myrpc/services/naming/proxy"
	"test/shared"
)

func main() {

	naming := namingproxy.New(shared.LocalHost, shared.NamingPort)
	calc := proxies.NewCalculadoraProxy(naming.Find("Calculadora"))

	scanner := bufio.NewScanner(os.Stdin)

	for {
		var a, b int

		fmt.Print("Type first number: ")
		scanner.Scan()
		a, _ = strconv.Atoi(scanner.Text())

		fmt.Print("Type second number: ")
		scanner.Scan()
		b, _ = strconv.Atoi(scanner.Text())

		res := calc.Som(a, b)
		fmt.Println("Response: ", res)

		fmt.Println()
	}
}
