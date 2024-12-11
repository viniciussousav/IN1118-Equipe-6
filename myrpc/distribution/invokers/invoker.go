package invokers

import (
	"log"
	"test/myrpc/app/businesses"
	"test/myrpc/distribution/core"
	"test/myrpc/distribution/interceptors"
	"test/myrpc/infrastructure"
	"test/shared"
)

type Invoker struct {
	ior               shared.IOR
	locationForwarder *interceptors.LocationForwarder
	localObjects      map[string]interface{}
}

func NewInvoker(host string, port int, locationForwarder *interceptors.LocationForwarder) Invoker {

	localObject := map[string]interface{}{
		"Calculadora": businesses.Calculadora{},
	}

	return Invoker{
		ior:               shared.IOR{Host: host, Port: port},
		locationForwarder: locationForwarder,
		localObjects:      localObject}
}

func (i Invoker) Invoke() {
	s := infrastructure.NewServerRequestHandler(i.ior.Host, i.ior.Port)
	m := core.Marshaller{}
	packet := core.Packet{}

	var rep int

	c := businesses.Calculadora{}

	for {
		// Invoke ServerRequestHandler
		b := s.Receive()

		// Unmarshall packet
		packet = m.Unmarshall(b)

		// Extract request from publisher
		r := core.ExtractRequest(packet)
		log.Print("Request received:", r)

		if r.ObjKey == "Calculadora" {

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
		}

		// Prepare reply
		var params []interface{}
		params = append(params, rep)

		// Create reply packet
		replyPacket := core.CreateReplyPacket(params, 200)

		// Marshall packet
		b = m.Marshall(replyPacket)

		// Send marshalled packet
		s.Send(b)
		log.Print("Response sent:", replyPacket.Bd.RepBody)

	}
}
