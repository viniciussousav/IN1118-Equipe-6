package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"test/myrpc/distribution/proxies"
	"test/shared"
)

func main() {

	calc := proxies.NewCalculadoraProxy(shared.IOR{Host: shared.LocalHost, Port: 8080})
	listenUserInput(calc)
}

func listenUserInput(calc proxies.CalculadoraProxy) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		var a, b int

		fmt.Print("Type first number: ")
		scanner.Scan()
		a, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("error reading input 'a': %s", err)
			continue
		}

		fmt.Print("Type second number: ")
		scanner.Scan()
		b, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("error reading input 'b': %s", err)
			continue
		}

		statusCode, res := calc.Som(a, b)

		if statusCode == 200 {
			fmt.Printf("StatusCode: %d, Result: %.2f \n", statusCode, res.Result[0].(float64))
		} else if statusCode == 301 {
			newLocationPort := int(res.Result[0].(float64))
			redirectProxy := proxies.NewCalculadoraProxy(shared.IOR{Host: shared.LocalHost, Port: newLocationPort})
			statusCode, res = redirectProxy.Som(a, b)
			fmt.Printf("Result from new location: %.2f \n", res.Result[0].(float64))
		} else {
			fmt.Printf("StatusCode: %d, Error: %s \n", statusCode, res.Result[0].(string))
		}
	}
}
