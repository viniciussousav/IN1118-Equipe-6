package main

import (
	"fmt"
	"test/myrpc/distribution/proxies"
	"test/shared"
	"time"
)

func main() {

	calc := proxies.NewCalculadoraProxy(shared.IOR{Host: shared.LocalHost, Port: 8080})
	executeSample(calc)
}

func executeSample(calc proxies.CalculadoraProxy) {

	a, b := 1, 2

	for i := 0; i < shared.SampleSize; i++ {

		t1 := time.Now()
		statusCode, res := calc.Som(a, b)

		switch statusCode {
		case 200:
			fmt.Printf("StatusCode: %d, Result: %.2f \n", statusCode, res.Result[0].(float64))
		case 301:
			newLocationPort := int(res.Result[0].(float64))
			redirectProxy := proxies.NewCalculadoraProxy(shared.IOR{Host: shared.LocalHost, Port: newLocationPort})
			statusCode, res = redirectProxy.Som(a, b)
			fmt.Printf("Result from new location: %.2f \n", res.Result[0].(float64))
		case 400:
			fmt.Printf("StatusCode: %d, Error: %s \n", statusCode, res.Result[0].(string))
		default:
			fmt.Printf("StatusCode: %d, Unknown threatment", statusCode)
		}

		fmt.Println(i, ";", time.Now().Sub(t1).Nanoseconds())
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Println("Experiment finalised...")
}
