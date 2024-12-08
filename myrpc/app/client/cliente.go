package main

import (
	"fmt"
	"test/myrpc/distribution/proxies/calculadora"
	namingproxy "test/myrpc/services/naming/proxy"
	"test/shared"
	"time"
)

func main() {
	Cliente()
}

func Cliente() {

	naming := namingproxy.New(shared.LocalHost, shared.NamingPort)
	calc := calculadoraproxy.NewCalculadoraProxy(naming.Find("Calculadora"))

	for i := 0; i < shared.StatisticSample; i++ {
		t1 := time.Now()
		for j := 0; j < shared.SampleSize; j++ {
			calc.Som(1, 2)
		}
		fmt.Println(i, ";", time.Now().Sub(t1).Milliseconds())
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Experiment finished...")
}
