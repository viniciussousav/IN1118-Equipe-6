package proxies

import (
	"test/location-forwarder-demo/distribution/requestor"
	"test/shared"
)

type CalculadoraProxy struct {
	Ior shared.IOR
}

func NewCalculadoraProxy(i shared.IOR) CalculadoraProxy {
	r := CalculadoraProxy{Ior: i}
	return r
}

func (p *CalculadoraProxy) Som(objectKey string, p1 int, p2 int) (statusCode int, content shared.Reply) {

	params := make([]interface{}, 2)
	params[0] = p1
	params[1] = p2

	req := shared.Request{ObjKey: objectKey, Op: "Som", Params: params}
	inv := shared.Invocation{Ior: p.Ior, Request: req}

	requestor := requestor.Requestor{}
	r := requestor.Invoke(inv)

	return r.Status, r.Rep
}
