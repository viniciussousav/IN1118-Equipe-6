package calculadorainvoker

import (
	"log"
	"math/rand"
	"test/myrpc/app/businesses/calculadora"
	"test/myrpc/distribution/interceptors/locationforwarder"
	"test/myrpc/distribution/marshaller"
	"test/myrpc/distribution/miop"
	"test/myrpc/infrastructure/srh"
	"test/shared"
)

type CalculadoraInvoker struct {
	Ior shared.IOR
}

func New(h string, p int) CalculadoraInvoker {
	ior := shared.IOR{Host: h, Port: p}
	inv := CalculadoraInvoker{Ior: ior}

	return inv
}

func (i CalculadoraInvoker) Invoke() {
	s := srh.NewSRH(i.Ior.Host, i.Ior.Port)
	m := marshaller.Marshaller{}
	miopPacket := miop.Packet{}
	var rep int

	locationForwarder := locationforwarder.LocationForwarder{}
	c := calculadora.Calculadora{}

	for {
		log.Print("Received!")

		// Invoke SRH
		b := s.Receive()

		isAvailable := rand.Intn(2) == 1

		if !isAvailable {
			locationForwarder.GetLocation("Calculator")
		}

		// Unmarshall miop packet
		miopPacket = m.Unmarshall(b)

		// Extract request from publisher
		r := miop.ExtractRequest(miopPacket)

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
		miop := miop.CreateReplyMIOP(params)

		// Marshall miop packet
		b = m.Marshall(miop)

		// Send marshalled packet
		s.Send(b)
	}
}
