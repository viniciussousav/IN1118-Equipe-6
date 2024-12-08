package calculadorainvoker

import (
	"log"
	"math/rand"
	"test/myrpc/app/businesses/calculadora"
	"test/myrpc/distribution/core"
	"test/myrpc/distribution/interceptors"
	"test/myrpc/infrastructure"
	"test/shared"
)

type Invoker struct {
	Ior shared.IOR
}

func NewInvoker(h string, p int) Invoker {
	ior := shared.IOR{Host: h, Port: p}
	inv := Invoker{Ior: ior}

	return inv
}

func (i Invoker) Invoke() {
	s := infrastructure.NewServerRequestHandler(i.Ior.Host, i.Ior.Port)
	m := core.Marshaller{}
	miopPacket := core.Packet{}
	var rep int

	locationForwarder := interceptors.LocationForwarder{}
	c := calculadora.Calculadora{}

	for {
		log.Print("Received!")

		// Invoke ServerRequestHandler
		b := s.Receive()

		isAvailable := rand.Intn(2) == 1

		if !isAvailable {
			locationForwarder.GetLocation("Calculator")
		}

		// Unmarshall miop packet
		miopPacket = m.Unmarshall(b)

		// Extract request from publisher
		r := core.ExtractRequest(miopPacket)

		_p1 := int(r.Params[0].(float64))
		_p2 := int(r.Params[1].(float64))

		// Demultiplex request & invoke Location Forwarder
		switch r.Op {
		case "Som":
			rep = c.Som(_p1, _p2)
		case "Dif":
			rep = c.Dif(_p1, _p2)
		case "Mul":
			rep = c.Mul(_p1, _p2)
		case "Div":
			rep = c.Div(_p1, _p2)
		default:
			log.Fatal("Invoker:: Operation '" + r.Op + "' is unknown:: ")
		}

		// Prepare reply
		var params []interface{}
		params = append(params, rep)

		// Create miop reply packet
		miop := core.CreateReplyMIOP(params)

		// Marshall miop packet
		b = m.Marshall(miop)

		// Send marshalled packet
		s.Send(b)
	}
}
