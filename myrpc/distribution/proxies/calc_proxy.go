package proxies

import (
	"log"
	"test/myrpc/distribution/requestor"
	"test/shared"
)

type CalculadoraProxy struct {
	Ior shared.IOR
}

func NewCalculadoraProxy(i shared.IOR) CalculadoraProxy {
	r := CalculadoraProxy{Ior: i}
	return r
}

func (p *CalculadoraProxy) Som(p1, p2 int) (statusCode int, content shared.Reply) {

	// 1. Configure input parameters
	params := make([]interface{}, 2)
	params[0] = p1
	params[1] = p2

	// Configure remote request
	req := shared.Request{ObjKey: "Calculadora", Op: "Som", Params: params}

	// Prepare invocation to Requestor
	inv := shared.Invocation{Ior: p.Ior, Request: req}

	// 3. Invoke Requestor
	requestor := requestor.Requestor{}
	r := requestor.Invoke(inv)

	log.Print(r)
	//4. Return something to the publisher
	return r.Status, r.Rep
}

func (p *CalculadoraProxy) Dif(p1, p2 int) int {

	// 1. Configure input parameters
	params := make([]interface{}, 2)
	params[0] = p1
	params[1] = p2

	// Configure remote request
	req := shared.Request{ObjKey: "Calculadora", Op: "Dif", Params: params}

	// Prepare invocation to Requestor
	inv := shared.Invocation{Ior: p.Ior, Request: req}

	// 3. Invoke Requestor
	requestor := requestor.Requestor{}
	r := requestor.Invoke(inv)

	//4. Return something to the publisher
	return int(r.Rep.Result[0].(float64)) // TODO
}
