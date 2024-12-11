package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"test/location-forwarder-demo/distribution/proxies"
	"test/shared"
	"time"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	calc := proxies.NewCalculadoraProxy(shared.IOR{Host: shared.LocalHost, Port: 8080})

	for {
		fmt.Println("New request...")

		fmt.Print("Object Key: ")
		scanner.Scan()
		objectKey := scanner.Text()

		fmt.Print("First number: ")
		scanner.Scan()
		a, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Invalid number! Try again.")
			continue
		}

		fmt.Print("Second number: ")
		scanner.Scan()
		b, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Invalid number! Try again.")
			continue
		}

		statusCode, res := calc.Som(objectKey, a, b)

		switch statusCode {
		case 200:
			log.Printf("StatusCode: %d, Result: %.2f", statusCode, res.Result[0].(float64))
		case 301:
			newLocationPort := int(res.Result[0].(float64))
			redirectProxy := proxies.NewCalculadoraProxy(shared.IOR{Host: shared.LocalHost, Port: newLocationPort})
			statusCode, res = redirectProxy.Som(objectKey, a, b)
			log.Printf("Result from new location: %.2f", res.Result[0].(float64))
		case 404:
			log.Printf("StatusCode: %d, Error: %s", statusCode, res.Result[0].(string))
		default:
			log.Printf("StatusCode: %d, Unknown threatment", statusCode)
		}
	}

	executeSample(calc)
}

func executeSample(calc proxies.CalculadoraProxy) {

	a, b := 1, 2
	objectKey := "Calculadora"

	for i := 0; i < shared.SampleSize; i++ {

		t1 := time.Now()
		statusCode, res := calc.Som(objectKey, a, b)

		switch statusCode {
		case 200:
			log.Printf("StatusCode: %d, Result: %.2f", statusCode, res.Result[0].(float64))
		case 301:
			newLocationPort := int(res.Result[0].(float64))
			redirectProxy := proxies.NewCalculadoraProxy(shared.IOR{Host: shared.LocalHost, Port: newLocationPort})
			statusCode, res = redirectProxy.Som(objectKey, a, b)
			log.Printf("Result from new location: %.2f", res.Result[0].(float64))
		case 400:
			log.Printf("StatusCode: %d, Error: %s", statusCode, res.Result[0].(string))
		default:
			log.Printf("StatusCode: %d, Unknown threatment", statusCode)
		}

		log.Println(i, ";", time.Now().Sub(t1).Nanoseconds())
		time.Sleep(time.Millisecond * 100)
	}

	log.Println("Experiment finalised...")
}
